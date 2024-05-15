package handle

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"codeberg.org/snonux/gos/internal"
	"codeberg.org/snonux/gos/internal/config/server"
	"codeberg.org/snonux/gos/internal/easyhttp"
	"codeberg.org/snonux/gos/internal/server/repository"
	"codeberg.org/snonux/gos/internal/types"
)

var getIDRe = regexp.MustCompile(`^/[0-9]{4}/[a-z0-9]{64}\.json$`)

// TODO: Use repository.Repository to store the file to the file system
func Submit(w http.ResponseWriter, r *http.Request, dataDir string) error {
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
	filePath := fmt.Sprintf("%s/%s/%s.json", dataDir, time.Now().Format("2006"), entry.ID)

	jsonStr, err := entry.Serialize()
	if err != nil {
		return err
	}

	if err := internal.SaveFile(filePath, jsonStr); err != nil {
		return err
	}

	return nil
}

func List(w http.ResponseWriter, r *http.Request, dataDir string) error {
	if r.Method != "GET" {
		return fmt.Errorf("expexted GET request")
	}

	list, err := repository.Instance(dataDir).List()
	if err != nil {
		return err
	}

	_, err = w.Write(list)
	return err
}

func Get(w http.ResponseWriter, r *http.Request, dataDir string) error {
	id := r.URL.Query().Get("id")
	if !getIDRe.MatchString(id) {
		return fmt.Errorf("invalid id %s", id)
	}

	data, err := os.ReadFile(fmt.Sprintf("%s/%s", dataDir, id))
	if err != err {
		return err
	}

	fmt.Fprint(w, string(data))
	return nil
}

func Merge(w http.ResponseWriter, r *http.Request, conf server.ServerConfig) error {
	var errs []error

	for _, partner := range conf.Partners() {
		if err := mergeFromPartner(conf, partner); err != nil {
			errs = append(errs, err)
		}
	}

	err := combineErrors(errs)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return err
	}

	fmt.Fprint(w, "Okiedokie")
	return nil
}

func mergeFromPartner(conf server.ServerConfig, partner string) error {
	var (
		errs  []error
		uri   = fmt.Sprintf("%s/list", partner)
		repo  = repository.Instance(conf.DataDir)
		pairs []repository.EntryPair
	)

	if err := easyhttp.GetData(uri, conf.ApiKey, &pairs); err != nil {
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

		if err := easyhttp.GetData(uri, conf.ApiKey, &entry); err != nil {
			errs = append(errs, err)
			continue
		}

		// In theory, this should never happen
		if pair.ID != entry.ID {
			errs = append(errs, fmt.Errorf("pair ID %s does not match entry id %s", pair.ID, entry.ID))
			continue
		}

		repo.Merge(entry)
	}

	return combineErrors(errs)
}

func combineErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	var sb strings.Builder
	for i, err := range errs {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(err.Error())
	}

	return errors.New(sb.String())
}
