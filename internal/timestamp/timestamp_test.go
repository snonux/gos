package timestamp

import (
	"strings"
	"testing"
	"time"
)

func TestUpdateInFilename(t *testing.T) {
	var (
		filePath        = "gosdir/db/platforms/mastodon/1728240487.txt.20241009-232530.queued"
		nowTime         = time.Now()
		updatedFilePath = UpdateInFilename(filePath, -2)
		parts           = strings.Split(updatedFilePath, ".")
	)

	updatedTime, err := Parse(parts[len(parts)-2])
	if err != nil {
		t.Error(err)
	}

	t.Log(filePath, updatedFilePath, nowTime.Sub(updatedTime))
}
