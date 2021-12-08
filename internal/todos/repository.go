package todos

import (
	"context"
	"github.com/tvitcom/czthree/internal/entity"
	"github.com/tvitcom/czthree/internal/dto"
	"github.com/tvitcom/czthree/pkg/dbcontext"
	"github.com/tvitcom/czthree/pkg/log"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// repository persists todoss in database
type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

// Repository encapsulates the logic to access todoss from the data source.
type Repository interface {
	// Get returns the users with follow by specified limit-offset params.
	GetUsersWithLimitOffset(ctx context.Context, limit, offset int64) ([]entity.User, error)
	// Get returns the todos with the specified todos ID.
	GettodosById(ctx context.Context, id int64) (entity.todos, error)
	// Get returns the todos with the specified todos ID.
	GettodossLast(ctx context.Context) ([]entity.todos, error)
	// Get returns the Messages by specified User_ID.
	GetMessagesSendersByUserId(ctx context.Context, uid int64) ([]dto.MessageSender, error)
	// Get returns the todos by word filtered like function.
	GettodossSearch(ctx context.Context, word string) ([]entity.todos, error)
	// Get returns the todoss by user_id.
	GettodossByUserId(ctx context.Context, uid int64) ([]entity.todos, error)
	// Get returns the user with the specified user_id.
	GetUserById(ctx context.Context, id int64) (entity.User, error)
	// Get returns the user with the specified todos_id.
	GetUserBytodosId(ctx context.Context, aid int64) (entity.User, error)
	// Get returns the todos with the specified email.
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	// Count returns the number of todoss.
	Count(ctx context.Context) (int, error)
	// Query returns the list of todoss with the given offset and limit.
	Querytodoss(ctx context.Context, offset, limit int) ([]entity.todos, error)
	// Create saves a new user in the storage.
	CreateUser(ctx context.Context, user entity.User) (int64, error)
	// Update user in the storage.
	UpdateUser(ctx context.Context, user entity.User, uid int64) error
	// Update User lastlogin field
	UpdateUserLastlogin(ctx context.Context, uid int64, dtstring string) error
	// Update User todoss Picture
	UpdatetodossPicture(ctx context.Context, aid int64, field, fname string) error
	// Create saves a new todos in the storage.
	Createtodos(ctx context.Context, todos entity.todos) (int64, error)
	// Getting todos by user_id.
	GettodossDisplayByUserId(ctx context.Context, id int64) ([]dto.todosDisplay, error)
	// Create saves a new todos in the storage.
	CreateMessage(ctx context.Context, message entity.Message) error
	// Update updates the todos with given ID in the storage.
	Update(ctx context.Context, todos entity.todos) error
	// Delete removes the todos with given ID from the storage.
	DeletetodosById(ctx context.Context, id int64) error
	// GetCategory(ctx) returns the list of todoss with the given offset and limit.
	GetCategory(ctx context.Context) ([]entity.Category, error)
	// GetCategoryPath(ctx) returns the list categories by special tree pathes query.
	GetCategoryPath(ctx context.Context) ([]entity.CategoryPath, error)
}

// NewRepository creates a new todos repository
func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

// Get the todos with the specified ID from the database.
func (r repository) GettodosById(ctx context.Context, id int64) (entity.todos, error) {
	var todos entity.todos
	err := r.db.With(ctx).Select().Model(id, &todos)
	return todos, err
}

// Get the user with the specified ID from the database.
func (r repository) GetUserById(ctx context.Context, id int64) (entity.User, error) {
	var user entity.User
	err := r.db.With(ctx).Select().From("user").Where(dbx.HashExp{"user_id": id}).One(&user)
	return user, err
}

// Get the user with the specified todos_ID from the database.
func (r repository) GetUserBytodosId(ctx context.Context, aid int64) (entity.User, error) {
	var user entity.User
	sql := `select u.* from user u left join todos a on u.user_id=a.user_id where a.todos_id={:aid}`
	err := r.db.With(ctx).NewQuery(sql).Bind(dbx.Params{"aid": aid}).One(&user)
	return user, err
}

// Get the user with the specified ID from the database.
func (r repository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	var user entity.User
	err := r.db.With(ctx).Select().From("user").Where(dbx.HashExp{"email": email}).One(&user)
	return user, err
}

// Create saves a new todos record in the database.
// It returns the ID of the newly inserted todos record.
func (r repository) CreateUser(ctx context.Context, user entity.User) (int64, error) {
	err := r.db.With(ctx).Model(&user).Insert()
	if err != nil {
		return 0, err
	}
	return user.UserId, nil
}

// Create saves a new todos record in the database.
// It returns the ID of the newly inserted todos record.
func (r repository) Createtodos(ctx context.Context, todos entity.todos) (int64, error) {
	err := r.db.With(ctx).Model(&todos).Insert()
	if err != nil {
		return 0, err
	}
	return todos.todosId, nil
}

// return todos records in the database.
func (r repository) GettodossDisplayByUserId(ctx context.Context, uid int64) ([]dto.todosDisplay, error) {
	var todoss []dto.todosDisplay
	sql := `
	SELECT a.todos_id, c.name category_name, a.title, a.price, a.currency, a.moderator_id, a.created, a.active 
	FROM todos a, category c 
	WHERE a.category_id = c.category_id AND a.user_id = {:uid}
	GROUP BY a.category_id 
	ORDER BY a.created ASC 
	LIMIT 100;`
	err := r.db.With(ctx).NewQuery(sql).Bind(dbx.Params{"uid": uid}).All(&todoss)
	return todoss, err
}

// returns recently added todos records
func (r repository) GetUsersWithLimitOffset(ctx context.Context, limit, offset int64) ([]entity.User, error) {
	var items []entity.User
	err := r.db.With(ctx).
		Select().
		From("user").
		// Where(dbx.HashExp{"active": 0}).
		Offset(offset).
		Limit(limit).
		OrderBy("lastlogin").
		All(&items)
	return items, err
}

// Create saves a new todos record in the database.
// It returns the ID of the newly inserted todos record.
// 
func (r repository) CreateMessage(ctx context.Context, message entity.Message) error {
	return r.db.With(ctx).Model(&message).Insert()
}

// Update saves the changes to an user in the database.
func (r repository) UpdateUser(ctx context.Context, user entity.User, uid int64) error {
	dbxvar := dbx.Params{
			"name": user.Name,
			"tel": user.Tel,
		}
	if user.Passhash != "" {
		dbxvar["passhash"]= user.Passhash
	}
	if user.Picture != "" {
		dbxvar["picture"] = user.Picture
	}
	// UPDATE `users` SET `status`={:p0} WHERE `id`={:p1}
	_, err := r.db.With(ctx).Update("user", dbxvar, dbx.HashExp{
		"user_id": uid,
	}).Execute()
	return err
}

// Update saves the changes to an user in the database.
func (r repository) UpdateUserLastlogin(ctx context.Context, uid int64, dtstring string) error {
	dbxvar := dbx.Params{
			"lastlogin": dtstring,
		}
	// UPDATE `users` SET `status`={:p0} WHERE `id`={:p1}
	_, err := r.db.With(ctx).Update("user", dbxvar, dbx.HashExp{
		"user_id": uid,
	}).Execute()
	return err
}
// UpdatetodossPicture(ctx context.Context, aid int64, field, fname string) error
func (r repository) UpdatetodossPicture(ctx context.Context, aid int64, field, fname string) error {
	dbxvar := dbx.Params{
			field: fname,
		}
	_, err := r.db.With(ctx).Update("todos", dbxvar, dbx.HashExp{
		"todos_id": aid,
	}).Execute()
	return err
}

// Update saves the changes to an todos in the database.
func (r repository) Update(ctx context.Context, todos entity.todos) error {
	return r.db.With(ctx).Model(&todos).Update()
}

// Delete deletes an todos with the specified ID from the database.
func (r repository) DeletetodosById(ctx context.Context, id int64) error {
	_, err := r.db.With(ctx).Delete("todos", dbx.HashExp{"todos_id": id}).Execute()
	return err
}

// Count returns the number of the todos records in the database.
func (r repository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.With(ctx).Select("COUNT(*)").From("todos").Row(&count)
	return count, err
}

// Query retrieves the todos records with the specified offset and limit from the database.
func (r repository) Querytodoss(ctx context.Context, offset, limit int) ([]entity.todos, error) {
	var todoss []entity.todos
	err := r.db.With(ctx).
		Select().
		OrderBy("id").
		Offset(int64(offset)).
		Limit(int64(limit)).
		All(&todoss)
	return todoss, err
}

// Count returns the number of the todos records in the database.
func (r repository) GetCategory(ctx context.Context) ([]entity.Category, error) {
	var categories []entity.Category
	err := r.db.With(ctx).
		Select().
		From("category").
		OrderBy("category_id").
		All(&categories)
	return categories, err
}

// Count returns the number of the todos records in the database.
func (r repository) GetCategoryPath(ctx context.Context) ([]entity.CategoryPath, error) {
	var items []entity.CategoryPath
	sql := `
	select c.category_id, concat_ws("-",( 
		select parc.name  
		from category parc  
		where c.parent_category_id = parc.category_id 
	),c.name) as path  
	from category c  where c.parent_category_id > 0 
	UNION 
	select c.category_id, c.name 
	from category c 
	where c.parent_category_id = 0 and c.category_id not in (
		select cc.parent_category_id 
		from category cc
	) order by 1 asc;`
	err := r.db.With(ctx).NewQuery(sql).All(&items)
	return items, err
}

// returns recently added todos records
func (r repository) GettodossLast(ctx context.Context) ([]entity.todos, error) {
	var limit int = 99
	var items []entity.todos
	err := r.db.With(ctx).
		Select().
		From("todos").
		Where(dbx.HashExp{"active": 1}).
		// Offset(int64(0)).
		Limit(int64(limit)).
		OrderBy("created").
		All(&items)
	return items, err
}

func (r repository) GetMessagesSendersByUserId(ctx context.Context, uid int64) ([]dto.MessageSender, error) {
	var items []dto.MessageSender
	sql := `
	select m.sender_id, u.name, u.email, u.tel, m.sended, m.moderator_id 
	from message m left join user u on m.sender_id = u.user_id 
	where m.receiver_id = {:uid}
	order by m.sended ASC
	`
	err := r.db.With(ctx).NewQuery(sql).Bind(dbx.Params{"uid": uid}).All(&items)
	return items, err
}

func (r repository) GettodossSearch(ctx context.Context, word string) ([]entity.todos, error) {
	var limit int = 200
	var items []entity.todos
	err := r.db.With(ctx).
		Select().
		From("todos").
		Where(dbx.And(dbx.HashExp{"active": 1}, dbx.Like("nanopost", word))).
		// Offset(int64(0)).
		Limit(int64(limit)).
		OrderBy("created").
		All(&items)
	return items, err
}

func (r repository) GettodossByUserId(ctx context.Context, uid int64) ([]entity.todos, error) {
	var items []entity.todos
	err := r.db.With(ctx).
		Select().
		From("todos").
		Where(dbx.HashExp{"user_id": uid}).
		OrderBy("created").
		All(&items)
	return items, err
}