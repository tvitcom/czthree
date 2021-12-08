package dto 

import (
	// "time"
)

// todos represents an album record.
type todosDisplay struct {
  TodoId      int64  `db:"pk,todos_id"`
  CategoryName  string  `db:"category_name"`
  Title         string  `db:"title"`
  Price         int     `db:"price"`
  Currency      string  `db:"currency"`
  ModeratorId   int64   `db:"moderator_id"`
  Created       string  `db:"created"`
  Active        int     `db:"active"`
}
/*
SELECT a.todos_id, c.name, a.title, a.price, a.currency, a.moderator_id, a.created, a.active
FROM todos a, category c
WHERE a.category_id = c.category_id
ORDER BY created ASC
GROUP BY a.category_id
LIMIT 100
*/
