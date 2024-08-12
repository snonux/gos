package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"codeberg.org/snonux/gos/internal/config/server"
	"codeberg.org/snonux/gos/internal/easyhttp"
	"codeberg.org/snonux/gos/internal/server/repository"
	"codeberg.org/snonux/gos/internal/types"
)

type Handler struct {
	conf    server.ServerConfig
	getIdRe *regexp.Regexp
}

func New(conf server.ServerConfig) Handler {
	return Handler{
		conf:    conf,
		getIdRe: regexp.MustCompile(`^/[0-9]{4}/[a-z0-9]{64}\.json$`),
	}
}

func (h Handler) Submit(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("expected POST request, but got %s", r.Method)
	}

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	ent, err := types.NewEntry(bytes)
	if err != nil {
		return err
	}
	return repository.Instance(h.conf).Merge(ent)
}

func (h Handler) List(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return fmt.Errorf("expexted GET request")
	}

	list, err := repository.Instance(h.conf).ListBytes()
	if err != nil {
		return err
	}

	_, err = w.Write(list)
	return err
}

func (h Handler) Get(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	if !h.getIdRe.MatchString(id) {
		return fmt.Errorf("invalid id %s", id)
	}

	ent, ok := repository.Instance(h.conf).Get(id)
	if !ok {
		return fmt.Errorf("no entry with id %s found", id)
	}

	fmt.Fprint(w, ent.String())
	return nil
}

func (h Handler) Merge(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var errs []error

	for _, partner := range h.conf.Partners() {
		if err := h.mergeFromPartner(ctx, partner); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	fmt.Fprint(w, "Okiedokie")
	return nil
}

func (h Handler) mergeFromPartner(ctx context.Context, partner string) error {
	var (
		errs  []error
		uri   = fmt.Sprintf("%s/list", partner)
		repo  = repository.Instance(h.conf)
		pairs []repository.EntryPair
	)

	if err := easyhttp.GetData(ctx, uri, h.conf.APIKey, &pairs); err != nil {
		return err
	}

	for _, pair := range pairs {
		if repo.HasSameEntry(pair) {
			continue
		}

		var (
			ent types.Entry
			uri = fmt.Sprintf("%s/get?id=%s", partner, pair.ID)
		)

		if err := easyhttp.GetData(ctx, uri, h.conf.APIKey, &ent); err != nil {
			errs = append(errs, err)
			continue
		}

		// In theory, this should never happen
		if pair.ID != ent.ID {
			errs = append(errs, fmt.Errorf("pair ID %s does not match entry id %s", pair.ID, ent.ID))
			continue
		}

		errs = append(errs, repo.Merge(ent))
	}

	return errors.Join(errs...)
}
