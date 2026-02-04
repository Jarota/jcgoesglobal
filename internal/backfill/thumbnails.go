package backfill

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/jarota/jctravels/internal/storage"
)

func Thumbnails(imgDir string) *thumbnails {
	return &thumbnails{imgDir}
}

type thumbnails struct {
	imgDir string
}

func (t *thumbnails) Run() error {
	slog.Info("starting thumbnail backfill ...")

	entries, err := os.ReadDir(t.imgDir)
	if err != nil {
		return fmt.Errorf("failed to read dir %s: %w", t.imgDir, err)
	}

	for _, entry := range entries {
		// skip non-jpeg files
		isJPEG := strings.Contains(entry.Name(), ".jpg") ||
			strings.Contains(entry.Name(), ".jpeg")
		if !isJPEG {
			continue
		}

		if strings.Contains(entry.Name(), "thumbnail") ||
			strings.Contains(entry.Name(), "compressed") {
			// no need to create a thumbnail of a thumbnail
			continue
		}

		src := t.imgDir + entry.Name()
		err = storage.CreateThumbnail(src)
		if err != nil {
			return fmt.Errorf("failed to create thumbnail for %s: %w", src, err)
		}

		err = storage.CreateCompressed(src)
		if err != nil {
			return fmt.Errorf("failed to create compressed image for %s: %w", src, err)
		}
	}

	slog.Info("backfill complete :)")
	return nil
}
