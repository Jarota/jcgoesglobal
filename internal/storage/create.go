package storage

import (
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

const createPostSQL = `
	INSERT INTO posts (id, caption, author, hearts) VALUES ($1, $2, $3, $4);
`

func (s *store) CreatePost(postID, caption, author string, hearts int) error {
	res, err := s.db.Exec(createPostSQL, postID, caption, author, hearts)
	if err != nil {
		return fmt.Errorf("failed to exec create post sql: %w", err)
	}

	return handleInsertResult(res)
}

const createImageSQL = `
	INSERT INTO images (filename, post_id) VALUES ($1, $2);
`

func (s *store) CreateImages(postID string, fileHeaders []*multipart.FileHeader) error {
	for _, f := range fileHeaders {
		file, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open file header: %w", err)
		}
		defer file.Close()

		dst, err := os.Create(uploadDir + f.Filename)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", f.Filename, err)
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			return fmt.Errorf("failed to copy file: %w", err)
		}

		res, err := s.db.Exec(createImageSQL, f.Filename, postID)
		if err != nil {
			return fmt.Errorf("failed to exec create image sql: %w", err)
		}

		if err = handleInsertResult(res); err != nil {
			return err
		}
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
