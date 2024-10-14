package timestamp

import (
	"strings"
	"testing"
)

func TestUpdateInFilename(t *testing.T) {
	t.Parallel()

	var (
		filePath        = "gosdir/db/platforms/mastodon/1728240487.txt.20241009-232530.queued"
		nowTime         = NowTime()
		updatedFilePath = UpdateInFilename(filePath, -2)
		parts           = strings.Split(updatedFilePath, ".")
	)

	updatedTime, err := Parse(parts[len(parts)-2])
	if err != nil {
		t.Error(err)
	}
	if nowTime.Sub(updatedTime) != 0 {
		t.Error("expected no time difference here")
	}
}
