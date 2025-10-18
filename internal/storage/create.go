package storage

import (
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/h2non/bimg"
)

const thumbnailSize = 800

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

			err = CreateThumbnail(path)
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

func CreateThumbnail(path string) error {
	orig, err := bimg.Read(path)
	if err != nil {
		return fmt.Errorf("failed to read path %s: %w", path, err)
	}

	new, err := bimg.NewImage(orig).Thumbnail(thumbnailSize)
	if err != nil {
		return fmt.Errorf("failed to init image: %w", err)
	}

	dst := thumbnailPath(path)
	err = bimg.Write(thumbnailPath(path), new)
	if err != nil {
		return fmt.Errorf("failed to write path %s: %w", dst, err)
	}

	return nil
}

func thumbnailPath(path string) string {
	len := len(path)
	suffix := "-thumbnail"
	switch {
	case strings.Contains(path, ".jpg"):
		return path[:len-4] + suffix + ".jpg"
	case strings.Contains(path, ".jpeg"):
		return path[:len-5] + suffix + ".jpeg"
	default:
		return ""
	}
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
