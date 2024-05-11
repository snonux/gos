package handle

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"

	"codeberg.org/snonux/gos/internal"
	"codeberg.org/snonux/gos/internal/server/repository"
	"codeberg.org/snonux/gos/internal/types"
)

var getIDRe = regexp.MustCompile(`^/[0-9]{4}/[a-z0-9]{64}\.json$`)

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

	list, err := repository.New(dataDir).List()
	if err != nil {
		return err
	}

	_, err = w.Write(list)
	return err
}

func Get(w http.ResponseWriter, r *http.Request, dataDir string) error {
	path := r.URL.Query().Get("path")
	if !getIDRe.MatchString(path) {
		return fmt.Errorf("invalid path %s", path)
	}

	data, err := os.ReadFile(fmt.Sprintf("%s/%s", dataDir, path))
	if err != err {
		return err
	}

	fmt.Fprint(w, string(data))
	return nil
}

func Merge(w http.ResponseWriter, r *http.Request, dataDir string) error {

	return nil
}
