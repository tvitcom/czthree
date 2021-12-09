package Todo

import (
	"context"
	vld "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/tvitcom/czthree/internal/entity"
	"github.com/tvitcom/czthree/internal/dto"
	"github.com/tvitcom/czthree/pkg/log"
	"github.com/tvitcom/czthree/internal/config"
	"github.com/tvitcom/czthree/pkg/util" 
	// "regexp"
	"errors"
	"fmt"
)

// Agregator encapsulates usecase logic for albums
type Agregator interface {
	GetUsersWithLimitOffset(ctx context.Context, limit, offset int64) ([]User, error)
	GetUserById(ctx context.Context, id int64) (User, error)
	GetUserByTodoId(ctx context.Context, aid int64) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetTodoById(ctx context.Context, id int64) (Todo, error)
	GetTodoDisplayByUserId(ctx context.Context, uid int64) ([]TodoDisplay, error)
	GetTodoByUserId(ctx context.Context, uid int64) ([]Todo, error)
	DeleteTodoData(ctx context.Context, aid int64) error
	UpdateUser(ctx context.Context, fdata *ProfileForm,uid int64, avafile string) error
	CreateTodo(ctx context.Context, f *TodoForm, uid, perfid int64, dt string) (int64, error)
}

type agregator struct {
	repo   Repository
	logger log.Logger
}

type Todo struct {
	entity.Todo
}
type User struct {
	entity.User
}
type TodoDisplay struct {
	dto.TodoDisplay
}

type TodoForm struct {
  AuthorId    int64  `json:"author_id" form:"author_id"`
  PerfomerId  int64  `json:"perfomer_id" form:"perfomer_id"`
  Name   			string `json:"name" form:"name"`
  Status      int    `json:"status" form:"status"`
}

// Validate validates the TodoForm fields
// func (m TodoForm) Validate() error {
// 	return vld.ValidateStruct(&m,
// 		vld.Field(&m.AuthorId, vld.Required, vld.Length(2, 128), vld.Match(regexp.MustCompile(`^(([a-zA-Z' -]{2,128})|([а-яА-ЯЁёІіЇїҐґЄє' -]{2,128}))`))),
// 		vld.Field(&m.PerfomerId, vld.Required),//, is.Digit),
// 		vld.Field(&m.Name, vld.Required),//, is.Digit),
// 		vld.Field(&m.Status, vld.Required),//, is.Digit),
// 		)
// }

// LoginForm represents an album update request.
type LoginForm struct {
	Username          string  `json:"username"        form:"username"`
	CurrentPassword   string  `json:"current-password" form:"current-password"`
}
// ProfileForm represents an album update request.
type ProfileForm struct {
	Tel               string  `json:"tel"          form:"tel"`
	GivenName         string  `json:"given-name"   form:"given-name"`
	NewPassword       string  `json:"new-password" form:"new-password"`
	NewPasswordRepeat string  `json:"new-password-repeat" form:"new-password-repeat"`
}

type DeleteTodoForm struct {
  TodoId int64  `json:"Todo_id" form:"Todo_id"`
}

// Validate validates the LoginForm fields
func (m LoginForm) Validate() error {
	return vld.ValidateStruct(&m,
		vld.Field(&m.Username, vld.When(config.CFG.AppMode != "dev", vld.Required, is.Email).Else(vld.Required, is.EmailFormat)),
		vld.Field(&m.CurrentPassword, vld.Required, vld.Length(6, 128)),
	)
}

// Validate validates the ProfileForm fields
// func (m ProfileForm) Validate() error {
// 	return vld.ValidateStruct(&m,
// 		vld.Field(&m.Tel, vld.When(m.Tel != "", vld.Required, vld.Length(10, 21), vld.Match(regexp.MustCompile(`[\d\s\-\+\(\)]{10,21}`))).Else(vld.Nil)),
// 		vld.Field(&m.GivenName, vld.Required, vld.Length(2, 64), vld.Match(regexp.MustCompile(`^(([a-zA-Z' -]{2,128})|([а-яА-ЯЁёІіЇїҐґЄє' -]{2,128}))`))),
// 		vld.Field(&m.NewPassword, vld.When(m.NewPassword != "", vld.Length(6, 128), vld.Required)),
// 		vld.Field(&m.NewPasswordRepeat, vld.When(m.NewPassword != "", vld.By(passwordsEquals(m.NewPassword)))),
// 	)
// }

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

// Update user.
func (ag agregator) UpdateUser(ctx context.Context, pf *ProfileForm, uid int64, avafile string) error {
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
	    Role:     "10",
		}, 
		uid,
	)
}

func (ag agregator) GetUserById(ctx context.Context, id int64) (User, error) {
	if id == 0 {
		return User{}, nil
	}
	user, err := ag.repo.GetUserById(ctx, id)
	if err != nil {
		return User{}, err
	}
	return User{user}, nil
}

func (ag agregator) GetUserByTodoId(ctx context.Context, aid int64) (User, error) {
	if aid == 0 {
		return User{}, nil
	}
	user, err := ag.repo.GetUserByTodoId(ctx, aid)
	if err != nil {
		return User{}, err
	}
	return User{user}, nil
}

func (ag agregator) GetTodoById(ctx context.Context, id int64) (Todo, error) {
	if id == 0 {
		return Todo{}, nil
	}
	todo, err := ag.repo.GetTodoById(ctx, id)
	if err != nil {
		return Todo{}, err
	}
	return Todo{todo}, nil
}

func (ag agregator) GetTodoDisplayByUserId(ctx context.Context, uid int64) ([]TodoDisplay, error) {
	result := []TodoDisplay{}
	items, err := ag.repo.GetTodoDisplayByUserId(ctx, uid)
	if err != nil {
		return result, err
	}
	for _, item := range items {
		result = append(result, TodoDisplay{item})
	}
	return result, nil
}

func (ag agregator) GetUserByEmail(ctx context.Context, email string) (User, error) {
	user, err := ag.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return User{}, err
	}
	return User{user}, nil
}

func (ag agregator) GetUsersWithLimitOffset(ctx context.Context, limit, offset int64) ([]User, error) {
	items, err := ag.repo.GetUsersWithLimitOffset(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	result := []User{}
	for _, item := range items {
		result = append(result, User{item})
	}
	return result, nil
}

func (ag agregator) GetTodoByUserId(ctx context.Context, uid int64) ([]Todo, error) {
	items, err := ag.repo.GetTodoByUserId(ctx, uid)
	if err != nil {
		return nil, err
	}
	result := []Todo{}
	for _, item := range items {
		result = append(result, Todo{item})
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
