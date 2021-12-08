package entity

import (
	// "time"
)

// todos represents an album record.
type Todo struct {
  TodoId    int64  `db:"pk,todos_id"`
  AuthorId      int64  `db:"user_id"`
  PerfomerId  int64  `db:"category_id"`
  Name       string  `db:"title"`
  Status      int     `db:"status"`
}
