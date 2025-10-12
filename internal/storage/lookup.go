package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jarota/jctravels/internal/model"
)

const lookupAllSQL = `
	SELECT posts.id, posts.caption, posts.author, posts.date, images.filename, posts.created_at
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
		var filename sql.NullString
		var date, createdAt time.Time
		if err := rows.Scan(&id, &caption, &author, &date, &filename, &createdAt); err != nil {
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

		if filename.Valid {
			// Only need to append `filename` to `uploadDir` when returning
			// as the frontend is served from *within* `s.siteRoot`
			posts[id].Pics = append(posts[id].Pics, s.uploadDir+filename.String)
		}
	}

	res := make([]*model.Post, 0, len(posts))
	for _, p := range posts {
		res = append(res, p)
	}

	return res, nil
}
