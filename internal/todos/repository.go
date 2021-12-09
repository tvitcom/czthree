package Todo

import (
	"context"
	"github.com/tvitcom/czthree/internal/entity"
	"github.com/tvitcom/czthree/internal/dto"
	"github.com/tvitcom/czthree/pkg/dbcontext"
	"github.com/tvitcom/czthree/pkg/log"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// repository persists Todo in database
type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

// Repository encapsulates the logic to access Todo from the data source.
type Repository interface {
	// Get returns the users with follow by specified limit-offset params.
	GetUsersWithLimitOffset(ctx context.Context, limit, offset int64) ([]entity.User, error)
	// Get returns the Todo with the specified Todo ID.
	GetTodoById(ctx context.Context, id int64) (entity.Todo, error)
	// Get returns the Todo with the specified Todo ID.
	GetTodoLast(ctx context.Context) ([]entity.Todo, error)
	// Get returns the Todo by word filtered like function.
	GetTodoSearch(ctx context.Context, word string) ([]entity.Todo, error)
	// Get returns the Todo by user_id.
	GetTodoByUserId(ctx context.Context, uid int64) ([]entity.Todo, error)
	// Get returns the user with the specified user_id.
	GetUserById(ctx context.Context, id int64) (entity.User, error)
	// Get returns the user with the specified Todo_id.
	GetUserByTodoId(ctx context.Context, aid int64) (entity.User, error)
	// Get returns the Todo with the specified email.
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	// Count returns the number of Todo.
	Count(ctx context.Context) (int, error)
	// Query returns the list of Todo with the given offset and limit.
	QueryTodo(ctx context.Context, offset, limit int) ([]entity.Todo, error)
	// Create saves a new user in the storage.
	CreateUser(ctx context.Context, user entity.User) (int64, error)
	// Update user in the storage.
	UpdateTodo(ctx context.Context, td entity.Todo, uid int64) error
	// Update user in the storage.
	UpdateTodoStatus(ctx context.Context, td entity.Todo) error
	// Update user in the storage.
	UpdateUser(ctx context.Context, user entity.User, uid int64) error
	// Create saves a new Todo in the storage.
	CreateTodo(ctx context.Context, Todo entity.Todo) (int64, error)
	// Getting Todo by user_id.
	GetTodoDisplayByUserId(ctx context.Context, id int64) ([]dto.TodoDisplay, error)
	// Update updates the Todo with given ID in the storage.
	Update(ctx context.Context, Todo entity.Todo) error
	// Delete removes the Todo with given ID from the storage.
	DeleteTodoById(ctx context.Context, id int64) error
}

// NewRepository creates a new Todo repository
func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

// Get the Todo with the specified ID from the database.
func (r repository) GetTodoById(ctx context.Context, id int64) (entity.Todo, error) {
	var Todo entity.Todo
	err := r.db.With(ctx).Select().Model(id, &Todo)
	return Todo, err
}

// Get the user with the specified ID from the database.
func (r repository) GetUserById(ctx context.Context, id int64) (entity.User, error) {
	var user entity.User
	err := r.db.With(ctx).Select().From("user").Where(dbx.HashExp{"user_id": id}).One(&user)
	return user, err
}

// Get the user with the specified Todo_ID from the database.
func (r repository) GetUserByTodoId(ctx context.Context, aid int64) (entity.User, error) {
	var user entity.User
	sql := `select u.* from user u left join Todo a on u.user_id=a.user_id where a.Todo_id={:aid}`
	err := r.db.With(ctx).NewQuery(sql).Bind(dbx.Params{"aid": aid}).One(&user)
	return user, err
}

// Get the user with the specified ID from the database.
func (r repository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	var user entity.User
	err := r.db.With(ctx).Select().From("user").Where(dbx.HashExp{"email": email}).One(&user)
	return user, err
}

// Create saves a new Todo record in the database.
// It returns the ID of the newly inserted Todo record.
func (r repository) CreateUser(ctx context.Context, user entity.User) (int64, error) {
	err := r.db.With(ctx).Model(&user).Insert()
	if err != nil {
		return 0, err
	}
	return user.UserId, nil
}

// Create saves a new Todo record in the database.
// It returns the ID of the newly inserted Todo record.
func (r repository) CreateTodo(ctx context.Context, Todo entity.Todo) (int64, error) {
	err := r.db.With(ctx).Model(&Todo).Insert()
	if err != nil {
		return 0, err
	}
	return Todo.TodoId, nil
}

// return Todo records in the database.
func (r repository) GetTodoDisplayByUserId(ctx context.Context, uid int64) ([]dto.TodoDisplay, error) {
	var Todos []dto.TodoDisplay
	sql := `SELECT t.todo_id, t.author_id, u.name manager, t.perfomer_id, t.name, t.status FROM todo t, user u WHERE t.perfomer_id = u.user_id AND t.perfomer_id = {:uid} ORDER BY t.status ASC`
	err := r.db.With(ctx).NewQuery(sql).Bind(dbx.Params{"uid": uid}).All(&Todos)
	return Todos, err
}

// returns recently added Todo records
func (r repository) GetUsersWithLimitOffset(ctx context.Context, limit, offset int64) ([]entity.User, error) {
	var items []entity.User
	err := r.db.With(ctx).
		Select().
		From("user").
		// Where(dbx.HashExp{"active": 0}).
		Offset(offset).
		Limit(limit).
		OrderBy("name").
		All(&items)
	return items, err
}

// Update saves the changes to an user in the database.
func (r repository) UpdateTodoStatus(ctx context.Context, td entity.Todo) error {
	dbxvar := dbx.Params{
			"status": td.Status,
		}
	// UPDATE `users` SET `status`={:p0} WHERE `id`={:p1}
	_, err := r.db.With(ctx).Update("todo", dbxvar, dbx.HashExp{
		"todo_id": td.TodoId,
	}).Execute()
	return err
}

// Update saves the changes to an user in the database.
func (r repository) UpdateTodo(ctx context.Context, td entity.Todo, uid int64) error {
	dbxvar := dbx.Params{
			"status": td.Status,
		}
	if td.Name != "" {
		dbxvar["name"] = td.Name
	}
	if td.PerfomerId != 0 {
		dbxvar["perfomer_id"] = td.PerfomerId
	}
	if td.PerfomerId != 0 {
		dbxvar["creator_id"] = td.PerfomerId
	}
	// UPDATE `users` SET `status`={:p0} WHERE `id`={:p1}
	_, err := r.db.With(ctx).Update("todo", dbxvar, dbx.HashExp{
		"todo_id": td.TodoId,
	}).Execute()
	return err
}


// Update saves the changes to an user in the database.
func (r repository) UpdateUser(ctx context.Context, user entity.User, uid int64) error {
	dbxvar := dbx.Params{
			"name": user.Name,
		}
	if user.Passhash != "" {
		dbxvar["passhash"]= user.Passhash
	}
	// UPDATE `users` SET `status`={:p0} WHERE `id`={:p1}
	_, err := r.db.With(ctx).Update("user", dbxvar, dbx.HashExp{
		"user_id": uid,
	}).Execute()
	return err
}

// Update saves the changes to an Todo in the database.
func (r repository) Update(ctx context.Context, Todo entity.Todo) error {
	return r.db.With(ctx).Model(&Todo).Update()
}

// Delete deletes an Todo with the specified ID from the database.
func (r repository) DeleteTodoById(ctx context.Context, id int64) error {
	_, err := r.db.With(ctx).Delete("Todo", dbx.HashExp{"Todo_id": id}).Execute()
	return err
}

// Count returns the number of the Todo records in the database.
func (r repository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.With(ctx).Select("COUNT(*)").From("Todo").Row(&count)
	return count, err
}

// Query retrieves the Todo records with the specified offset and limit from the database.
func (r repository) QueryTodo(ctx context.Context, offset, limit int) ([]entity.Todo, error) {
	var Todo []entity.Todo
	err := r.db.With(ctx).
		Select().
		OrderBy("id").
		Offset(int64(offset)).
		Limit(int64(limit)).
		All(&Todo)
	return Todo, err
}

// returns recently added Todo records
func (r repository) GetTodoLast(ctx context.Context) ([]entity.Todo, error) {
	var limit int = 99
	var items []entity.Todo
	err := r.db.With(ctx).
		Select().
		From("Todo").
		Where(dbx.HashExp{"active": 1}).
		// Offset(int64(0)).
		Limit(int64(limit)).
		OrderBy("created").
		All(&items)
	return items, err
}

func (r repository) GetTodoSearch(ctx context.Context, word string) ([]entity.Todo, error) {
	var limit int = 200
	var items []entity.Todo
	err := r.db.With(ctx).
		Select().
		From("Todo").
		Where(dbx.And(dbx.HashExp{"active": 1}, dbx.Like("nanopost", word))).
		// Offset(int64(0)).
		Limit(int64(limit)).
		OrderBy("created").
		All(&items)
	return items, err
}

func (r repository) GetTodoByUserId(ctx context.Context, uid int64) ([]entity.Todo, error) {
	var items []entity.Todo
	err := r.db.With(ctx).
		Select().
		From("Todo").
		Where(dbx.HashExp{"user_id": uid}).
		OrderBy("created").
		All(&items)
	return items, err
}