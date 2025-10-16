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
		filename := entry.Name()

		var newFilename string
		len := len(filename)
		suffix := "-thumbnail"
		switch {
		case strings.Contains(filename, "-thumbnail"):
			// no need to create thumbnails of thumbnails
			continue
		case strings.Contains(filename, ".jpg"):
			newFilename = filename[:len-4] + suffix + ".jpg"
		case strings.Contains(filename, ".jpeg"):
			newFilename = filename[:len-5] + suffix + ".jpeg"
		default:
			continue
		}

		// in the interest of idempotency, skip already created thumbnails
		_, err := os.Open(newFilename)
		if os.IsNotExist(err) {
			srcPath := t.imgDir + filename
			src, err := os.Open(srcPath)
			if err != nil {
				return fmt.Errorf("failed to open %s: %w", srcPath, err)
			}

			dst := t.imgDir + newFilename
			err = storage.CreateThumbnail(src, dst)
			if err != nil {
				return fmt.Errorf("failed to create thumbnail for %s: %w", dst, err)
			}
		}
	}

	slog.Info("backfill complete :)")
	return nil
}
