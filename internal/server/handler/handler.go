package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
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

// TODO: Use repository.Repository to store the file to the file system
func (h Handler) Submit(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("expexted POST request")
	}

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	entry, err := types.NewEntry(bytes)
	if err != nil {
		return err
	}
	return repository.Instance(h.conf.DataDir).Merge(entry)
}

func (h Handler) List(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return fmt.Errorf("expexted GET request")
	}

	list, err := repository.Instance(h.conf.DataDir).List()
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

	data, err := os.ReadFile(fmt.Sprintf("%s/%s", h.conf.DataDir, id))
	if err != err {
		return err
	}

	fmt.Fprint(w, string(data))
	return nil
}

func (h Handler) Merge(w http.ResponseWriter, r *http.Request) error {
	var errs []error

	for _, partner := range h.conf.Partners() {
		if err := h.mergeFromPartner(partner); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	fmt.Fprint(w, "Okiedokie")
	return nil
}

func (h Handler) mergeFromPartner(partner string) error {
	var (
		errs  []error
		uri   = fmt.Sprintf("%s/list", partner)
		repo  = repository.Instance(h.conf.DataDir)
		pairs []repository.EntryPair
	)

	if err := easyhttp.GetData(uri, h.conf.ApiKey, &pairs); err != nil {
		return err
	}

	for _, pair := range pairs {
		if repo.HasEntry(pair) {
			continue
		}

		var (
			entry types.Entry
			uri   = fmt.Sprintf("%s/get?id=%s", partner, pair.ID)
		)

		if err := easyhttp.GetData(uri, h.conf.ApiKey, &entry); err != nil {
			errs = append(errs, err)
			continue
		}

		// In theory, this should never happen
		if pair.ID != entry.ID {
			errs = append(errs, fmt.Errorf("pair ID %s does not match entry id %s", pair.ID, entry.ID))
			continue
		}

		errs = append(errs, repo.Merge(entry))
	}

	return errors.Join(errs...)
}
