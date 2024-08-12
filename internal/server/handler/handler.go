package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"codeberg.org/snonux/gos/internal/config/server"
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
		getIdRe: regexp.MustCompile(`^[a-z0-9]{64}$`),
		// getIdRe: regexp.MustCompile(`^/[0-9]{4}/[a-z0-9]{64}\.json$`),
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

	bytes, err := ent.JSONBytes()
	if err != nil {
		return err
	}

	fmt.Fprint(w, string(bytes))
	return nil
}

func (h Handler) Merge(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if err := repository.Instance(h.conf).MergeRemotely(ctx); err != nil {
		return err
	}

	fmt.Fprint(w, "Repository merge went well")
	return nil
}
