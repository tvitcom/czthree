package entity

// User represents a user.
type User struct {
  UserId       int64  `db:"user_id"`
  Name         string  `db:"name"`
  Email        string  `db:"email"`
  Passhash     string  `db:"passhash"`
  Role         int  `db:"role"`
}

// GetID returns the user ID.
func (u User) GetID() int64 {
	return u.UserId
}

// GetName returns the user name.
func (u User) GetName() string {
	return u.Name
}
