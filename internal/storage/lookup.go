package storage

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jarota/jctravels/internal/model"
)

const lookupAllSQL = `
	SELECT posts.id, posts.caption, posts.author, posts.date, images.id, images.filename, posts.created_at
	FROM posts LEFT JOIN images ON posts.id = images.post_id
	ORDER BY posts.created_at DESC;
`

func (s *store) LookupAll() ([]*model.Post, error) {
	rows, err := s.db.Query(lookupAllSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to query all posts: %w", err)
	}
	defer rows.Close()

	posts := make(map[string]*model.Post)
	for rows.Next() {
		var id, caption, author string
		var picID, filename sql.NullString
		var date, createdAt time.Time
		if err := rows.Scan(&id, &caption, &author, &date, &picID, &filename, &createdAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if _, ok := posts[id]; !ok {
			posts[id] = &model.Post{
				ID:        id,
				Caption:   caption,
				Author:    author,
				Date:      date,
				CreatedAt: createdAt,
			}
		}

		if picID.Valid && filename.Valid {
			var thumbnail string
			len := len(filename.String)
			suffix := "-thumbnail"
			switch {
			case strings.Contains(filename.String, ".jpg"):
				thumbnail = filename.String[:len-4] + suffix + ".jpg"
			case strings.Contains(filename.String, ".jpeg"):
				thumbnail = filename.String[:len-5] + suffix + ".jpeg"
			}

			var thumbnailPath string
			if thumbnail != "" {
				thumbnailPath = s.uploadDir + thumbnail
			}

			// Only need to append `filename` to `uploadDir` when returning
			// as the frontend is served from *within* `s.siteRoot`
			posts[id].Pics = append(
				posts[id].Pics,
				model.Pic{
					ID:            picID.String,
					HDPath:        s.uploadDir + filename.String,
					ThumbnailPath: thumbnailPath,
				},
			)
		}
	}

	res := make([]*model.Post, 0, len(posts))
	for _, p := range posts {
		res = append(res, p)
	}

	return res, nil
}
