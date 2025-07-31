package storage

import (
	"database/sql"
	"fmt"

	"github.com/jarota/jctravels/internal/model"
)

const lookupAllSQL = `
	SELECT posts.id, posts.caption, images.filename
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
		var id, caption string
		var filename sql.NullString
		if err := rows.Scan(&id, &caption, &filename); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if _, ok := posts[id]; !ok {
			posts[id] = &model.Post{
				ID:      id,
				Caption: caption,
			}
		}

		if filename.Valid {
			posts[id].Pics = append(posts[id].Pics, filename.String)
		}
	}

	res := make([]*model.Post, 0, len(posts))
	for _, p := range posts {
		res = append(res, p)
	}

	return res, nil
}
