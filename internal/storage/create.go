package storage

import (
	"database/sql"
	"fmt"
	"image/jpeg"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rymdport/resize"
)

const thumbnailSize = 500

const createPostSQL = `
	INSERT INTO posts (id, caption, author, date) VALUES ($1, $2, $3, $4);
`

func (s *store) CreatePost(postID, caption, author string, date time.Time) error {
	res, err := s.db.Exec(createPostSQL, postID, caption, author, date)
	if err != nil {
		return fmt.Errorf("failed to exec create post sql: %w", err)
	}

	return handleInsertResult(res)
}

const createImageSQL = `
	INSERT INTO images (id, filename, post_id) VALUES ($1, $2, $3);
`

func (s *store) CreateImages(postID string, fileHeaders []*multipart.FileHeader) error {
	for _, f := range fileHeaders {
		file, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open file header: %w", err)
		}
		defer file.Close()

		path := s.siteRoot + s.uploadDir + f.Filename
		dst, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", path, err)
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			return fmt.Errorf("failed to copy file: %w", err)
		}

		// only create thumbnails for jpeg files
		if strings.Contains(f.Filename, ".jpg") ||
			strings.Contains(f.Filename, ".jpeg") {

			err = CreateThumbnail(dst, path+"-thumbnail")
			if err != nil {
				return fmt.Errorf("failed to create thumbnail for %s: %w", path, err)
			}
		}

		res, err := s.db.Exec(createImageSQL, uuid.NewString(), f.Filename, postID)
		if err != nil {
			return fmt.Errorf("failed to exec create image sql: %w", err)
		}

		if err = handleInsertResult(res); err != nil {
			return err
		}
	}

	return nil
}

func CreateThumbnail(orig io.Reader, path string) error {
	img, err := jpeg.Decode(orig)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	dst, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", path, err)
	}
	defer dst.Close()

	// 0 height preserves aspect ratio
	thumbnail := resize.Resize(
		thumbnailSize,
		0,
		img,
		resize.MitchellNetravali,
	)
	err = jpeg.Encode(dst, thumbnail, nil)
	if err != nil {
		return fmt.Errorf("failed to encode thumbnail: %w", err)
	}

	return nil
}

func handleInsertResult(res sql.Result) error {
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if n != 1 {
		return fmt.Errorf("unexpected number of rows affected: %v", n)
	}
	return nil
}
