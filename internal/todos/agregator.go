package todos

import (
	"context"
	vld "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/tvitcom/czthree/internal/entity"
	"github.com/tvitcom/czthree/internal/dto"
	"github.com/tvitcom/czthree/pkg/log"
	"github.com/tvitcom/czthree/internal/config"
	"github.com/tvitcom/czthree/pkg/util" 
	"regexp"
	"errors"
	"fmt"
)

// Agregator encapsulates usecase logic for albums
type Agregator interface {
	GetUsersWithLimitOffset(ctx context.Context, limit, offset int64) ([]entity.User, error)
	GetUserById(ctx context.Context, id int64) (entity.User, error)
	GetUserByTodoId(ctx context.Context, aid int64) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	GetTodoById(ctx context.Context, id int64) (entity.Todo, error)
	GetTodoDisplayByUserId(ctx context.Context, uid int64) ([]dto.TodoDisplay, error)
	GetTodoByUserId(ctx context.Context, uid int64) ([]entity.Todo, error)
	UpdateTodoStatus(ctx context.Context, fdata *TodoStatusForm, uid int64) error
	UpdateTodo(ctx context.Context, fdata *TodoForm) error
	DeleteTodoData(ctx context.Context, aid int64) error
	UpdateUser(ctx context.Context, fdata *ProfileForm,uid int64) error
	CreateTodo(ctx context.Context, f *TodoForm, uid, perfid int64, dt string) (int64, error)
}

type agregator struct {
	repo   Repository
	logger log.Logger
}

// type Todo struct {
// 	entity.Todo
// }
// type User struct {
// 	entity.User
// }
// type TodoDisplay struct {
// 	dto.TodoDisplay
// }

type TodoForm struct {
  TodoId    int64  `json:"todo_id" form:"todo_id"`
  AuthorId    int64  `json:"author_id" form:"author_id"`
  PerfomerId  int64  `json:"perfomer_id" form:"perfomer_id"`
  Name   			string `json:"name" form:"name"`
  Status      int    `json:"status" form:"status"`
}

// Validate validates the TodoForm fields
func (m TodoForm) Validate() error {
	return vld.ValidateStruct(&m,
		vld.Field(&m.TodoId, vld.Required),
		vld.Field(&m.AuthorId, vld.Required),
		vld.Field(&m.PerfomerId, vld.Required),//, is.Digit),
		vld.Field(&m.Name, vld.Required, vld.Length(2, 128)),//, is.Digit),
		vld.Field(&m.Status, vld.Required),//, is.Digit),
		)
}

// LoginForm represents an album update request.
type LoginForm struct {
	Username          string  `json:"username"        form:"username"`
	CurrentPassword   string  `json:"current-password" form:"current-password"`
}
// pdateTodoForm represents an album update request.
type TodoStatusForm struct {
	TodoId       int64   `json:"todoid" form:"todoit"`
	Status       int     `json:"status" form:"status"`
}
// ProfileForm represents an album update request.
type ProfileForm struct {
	Tel               string  `json:"tel"          form:"tel"`
	GivenName         string  `json:"given-name"   form:"given-name"`
	NewPassword       string  `json:"new-password" form:"new-password"`
	NewPasswordRepeat string  `json:"new-password-repeat" form:"new-password-repeat"`
}

type DeleteTodoForm struct {
  TodoId int64  `json:"todo_id" form:"todo_id"`
}

// Validate validates the LoginForm fields
func (m LoginForm) Validate() error {
	return vld.ValidateStruct(&m,
		vld.Field(&m.Username, vld.When(config.CFG.AppMode != "dev", vld.Required, is.Email).Else(vld.Required, is.EmailFormat)),
		vld.Field(&m.CurrentPassword, vld.Required, vld.Length(6, 128)),
	)
}

// Validate validates the LoginForm fields
func (f TodoStatusForm) Validate() error {
	return vld.ValidateStruct(&f,
		vld.Field(&f.TodoId, vld.Required),
		vld.Field(&f.Status, vld.Required),
	)
}

// Validate validates the ProfileForm fields
func (m ProfileForm) Validate() error {
	return vld.ValidateStruct(&m,
		// vld.Field(&m.Tel, vld.When(m.Tel != "", vld.Required, vld.Length(10, 21), vld.Match(regexp.MustCompile(`[\d\s\-\+\(\)]{10,21}`))).Else(vld.Nil)),
		vld.Field(&m.GivenName, vld.Required, vld.Length(2, 64), vld.Match(regexp.MustCompile(`^(([a-zA-Z' -]{2,128})|([а-яА-ЯЁёІіЇїҐґЄє' -]{2,128}))`))),
		vld.Field(&m.NewPassword, vld.When(m.NewPassword != "", vld.Length(6, 128), vld.Required)),
		vld.Field(&m.NewPasswordRepeat, vld.When(m.NewPassword != "", vld.By(passwordsEquals(m.NewPassword)))),
	)
}

func passwordsEquals(str string) vld.RuleFunc {
	return func(value interface{}) error {
		s, _ := value.(string)
        if s != str {
            return errors.New("unexpected string")
        }
        return nil
    }
}

// NewAgregator creates a new album agregator.
func NewAgregator(repo Repository, logger log.Logger) agregator {
	return agregator{repo, logger}
}

// Update sattus in todo.
func (ag agregator) UpdateTodoStatus(ctx context.Context, tdf *TodoStatusForm, uid int64) error {
	var newstatus int
	if uid != config.AdminUserID && tdf.Status == config.TODO_DONE { // status == 2
		newstatus = config.TODO_DONE
	} else if uid == config.AdminUserID {
		newstatus = tdf.Status
	} else {
		return errors.New("Bad user credential for operation")
	}
	return ag.repo.UpdateTodoStatus(ctx, 
		entity.Todo{
	    TodoId: tdf.TodoId,
	    Status: newstatus,
		},
	)
}

// Update sattus in todo.
func (ag agregator) UpdateTodo(ctx context.Context, tdf *TodoForm) error {
	return ag.repo.UpdateTodo(ctx, 
		entity.Todo{
	    TodoId: tdf.TodoId,
	    AuthorId: tdf.AuthorId,
	    PerfomerId: tdf.PerfomerId,
	    Name: tdf.Name,
	    Status: tdf.Status,
		})
}

// Update user.
func (ag agregator) UpdateUser(ctx context.Context, pf *ProfileForm, uid int64) error {
  // Save new User record
  if pf.NewPassword != "" {
      pf.NewPassword = util.MakeBCryptHash(pf.NewPassword, config.BCRYPT_COST)
  } else {
      pf.NewPassword = ""
  }
	return ag.repo.UpdateUser(ctx, 
		entity.User{
	    Name:     pf.GivenName,
	    Passhash: pf.NewPassword,
	    Role:     config.USER_ROLE,
		}, 
		uid,
	)
}

func (ag agregator) GetUserById(ctx context.Context, id int64) (entity.User, error) {
	if id == 0 {
		return entity.User{}, nil
	}
	user, err := ag.repo.GetUserById(ctx, id)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (ag agregator) GetUserByTodoId(ctx context.Context, aid int64) (entity.User, error) {
	if aid == 0 {
		return entity.User{}, nil
	}
	user, err := ag.repo.GetUserByTodoId(ctx, aid)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (ag agregator) GetTodoById(ctx context.Context, id int64) (entity.Todo, error) {
	if id == 0 {
		return entity.Todo{}, nil
	}
	todo, err := ag.repo.GetTodoById(ctx, id)
fmt.Printf("AGREG-TODO: %v", todo)
	if err != nil {
		return entity.Todo{}, err
	}
	return todo, nil
}

func (ag agregator) GetTodoDisplayByUserId(ctx context.Context, uid int64) ([]dto.TodoDisplay, error) {
	result := []dto.TodoDisplay{}
	items, err := ag.repo.GetTodoDisplayByUserId(ctx, uid)
	if err != nil {
		return result, err
	}
	for _, item := range items {
		result = append(result, item)
	}
	return result, nil
}

func (ag agregator) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	user, err := ag.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (ag agregator) GetUsersWithLimitOffset(ctx context.Context, limit, offset int64) ([]entity.User, error) {
	items, err := ag.repo.GetUsersWithLimitOffset(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	result := []entity.User{}
	for _, item := range items {
		result = append(result, item)
	}
	return result, nil
}

func (ag agregator) GetTodoByUserId(ctx context.Context, uid int64) ([]entity.Todo, error) {
	items, err := ag.repo.GetTodoByUserId(ctx, uid)
	if err != nil {
		return nil, err
	}
	result := []entity.Todo{}
	for _, item := range items {
		result = append(result, item)
	}
	return result, nil
}

func (ag agregator) DeleteTodoData(ctx context.Context, aid int64) error {
	aidStr := fmt.Sprint(aid)
	pathes := config.PictureTodoPath + aidStr + "_*"
  err := ag.repo.DeleteTodoById(ctx, aid)
  if err != nil {
      ag.logger.With(ctx).Error(err.Error())
      return err
  }
	return util.FileDeletionByMask(ctx, pathes)
}

// Create creates a new Todo.
func (ag agregator) CreateTodo(ctx context.Context, f *TodoForm, uid, perfid int64, name string) (int64, error) {
	if uid == 0 {
		return 0, errors.New("CreateTodo on User_id:0")
	}
	return ag.repo.CreateTodo(ctx, entity.Todo{
	  AuthorId: uid,
	  PerfomerId: perfid,
	  Name: name,
	  Status: 1, // todo status:1
	})
}
