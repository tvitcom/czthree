package dto 

import (
	// "time"
)

// Todo represents an album record:t.todo_id, u.name manager, t.perfomer_id, t.name, t.status 
type TodoDisplay struct {
  TodoId        int64  `db:"pk,todo_id"`
  PerfomerId    int64  `db:"perfomer_id"`
  Manager       string  `db:"manager"`
  Name          string `db:"name"`
  Status        int    `db:"status"`
}
