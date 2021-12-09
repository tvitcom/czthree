package dto 

import (
	// "time"
)

// Todo represents an album record.
type TodoDisplay struct {
  TodoId        int64  `db:"pk,todo_id"`
  AuthorId      int64  `db:"author_id"`
  PerfomerId    int64  `db:"perfomer_id"`
  Name          string `db:"name"`
  Status        int    `db:"status"`
}
