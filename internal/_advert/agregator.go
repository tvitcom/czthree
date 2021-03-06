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
	// "github.com/tvitcom/czthree/pkg/dbcontext"
	"strings"
	"regexp"
	"errors"
	"time"
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
	GetTodoLast(ctx context.Context) ([]Todo, error)
	GetTodoSearch(ctx context.Context, clause string) ([]Todo, error)
	GetTodoByUserId(ctx context.Context, uid int64) ([]Todo, error)
	DeleteTodoData(ctx context.Context, aid int64) error
	CreateUser(ctx context.Context, fdata *SignupForm, dtstring, roles, notes, avafile string) (int64, error)
	UpdateUser(ctx context.Context, fdata *ProfileForm,uid int64, avafile string) error
	CreateTodo(ctx context.Context, f *QuickTodoForm, uid int64, dt string) (int64, error)
}

type agregator struct {
	repo   Repository
	logger log.Logger
}

type Todo struct {
	entity.Todo
}
type Category struct {
	entity.Category
}
type CategoryPath struct {
	entity.CategoryPath
}
type User struct {
	entity.User
}
type Message struct {
	entity.Message
}
type MessageSender struct {
	dto.MessageSender
}
type TodoDisplay struct {
	dto.TodoDisplay
}

type QuickTodoForm struct {
  Email       string  `json:"email" form:"email"`
  GivenName   string  `json:"given-name" form:"given-name"`
  CategoryId  int64   `json:"category_id" form:"category_id"`
  Price    		string  `json:"price" form:"price"`
  Nanopost    string  `json:"nanopost" form:"nanopost"`
  Tel         string  `json:"tel" form:"tel"`
}

type QuickSearchForm struct {
  CategoryId int64  `json:"category-id" form:"category-id"`
  Clause string     `json:"q" form:"q"`
}

// Validate validates the QuickTodoForm fields
func (m QuickTodoForm) Validate() error {
	return vld.ValidateStruct(&m,
		vld.Field(&m.GivenName, vld.Required, vld.Length(2, 128), vld.Match(regexp.MustCompile(`^(([a-zA-Z' -]{2,128})|([??-????-??????????????????????' -]{2,128}))`)), vld.By(swearDetector())),
		vld.Field(&m.Email, vld.When(config.CFG.AppMode != "dev", vld.Required, is.Email).Else(vld.Required, is.EmailFormat)),
		vld.Field(&m.CategoryId, vld.Required),//, is.Digit),
		vld.Field(&m.Price, vld.Length(0, 45)),
		vld.Field(&m.Nanopost, vld.Required, vld.Length(12, 512), vld.By(swearDetector())),
		vld.Field(&m.Tel, vld.When(m.Tel != "", vld.Required, vld.Length(10, 21), vld.Match(regexp.MustCompile(`^[\-\+\d\s\(\)]{10,21}$`))).Else(vld.Nil)),
	)
}

// LoginForm represents an album update request.
type LoginForm struct {
	Username          string  `json:"username"        form:"username"`
	CurrentPassword   string  `json:"current-password" form:"current-password"`
}

// SignupForm represents an album update request.
type SignupForm struct {
	Email             string  `json:"email"       form:"email"`
	GivenName         string  `json:"given-name"  form:"given-name"`
	NewPassword       string  `json:"new-password" form:"new-password"`
	NewPasswordRepeat string  `json:"new-password-repeat" form:"new-password-repeat"`
	Tel               string  `json:"tel"          form:"tel"`
}
// ProfileForm represents an album update request.
type ProfileForm struct {
	Tel               string  `json:"tel"          form:"tel"`
	GivenName         string  `json:"given-name"   form:"given-name"`
	NewPassword       string  `json:"new-password" form:"new-password"`
	NewPasswordRepeat string  `json:"new-password-repeat" form:"new-password-repeat"`
}
// MessageForm represents an album update request.
type MessageForm struct {
	TodoId    int64   `json:"Todo-id"        form:"Todo-id"`
	Email       string  `json:"email"        form:"email"`
	GivenName   string  `json:"given-name"   form:"given-name"`
	Msg         string  `json:"msg"          form:"msg"`
}

// SupportForm represents an album update request.
type SupportForm struct {
	Email       string  `json:"email"        form:"email"`
	Subject     string  `json:"subject"      form:"subject"`
	GivenName   string  `json:"given-name"   form:"given-name"`
	Msg         string  `json:"msg"          form:"msg"`
}

type DeleteTodoForm struct {
  TodoId int64  `json:"Todo_id" form:"Todo_id"`
}

type WatchAuthorForm struct {
  RecaptchaResponse string  `json:"g-recaptcha-response" form:"g-recaptcha-response"`
  TodoId int64  `json:"Todo_id" form:"Todo_id"`
  SignerUA string  `json:"signer_ua" form:"signer_ua"`
  SignerScreen string  `json:"signer_screen" form:"signer_screen"`
  SignerLangs string  `json:"signer_langs" form:"signer_langs"`
  SignerTime string  `json:"signer_time" form:"signer_time"`
}

// Validate validates the SignupForm fields
func (m SignupForm) Validate() error {
	return vld.ValidateStruct(&m,
		vld.Field(&m.Email, vld.When(config.CFG.AppMode != "dev", vld.Required, is.Email).Else(vld.Required, is.EmailFormat)),
		vld.Field(&m.GivenName, vld.Required, vld.Length(2, 64), vld.Match(regexp.MustCompile(`^(([a-zA-Z' -]{2,128})|([??-????-??????????????????????' -]{2,128}))`))),
		vld.Field(&m.NewPassword, vld.Required, vld.Length(6, 128)),
		vld.Field(&m.NewPasswordRepeat, vld.Required, vld.Length(6, 128)),
		vld.Field(&m.NewPasswordRepeat, vld.By(passwordsEquals(m.NewPassword))),
		vld.Field(&m.Tel, vld.When(m.Tel != "", vld.Required, vld.Length(10, 21), vld.Match(regexp.MustCompile(`^[\d\s\-\+\(\)]{10,21}$`))).Else(vld.Nil)),
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

func swearDetector() vld.RuleFunc {
	return func(value interface{}) error {
			s, _ := util.SwearDetector(value.(string))
	    if len(s) > 0 {
	        return errors.New("???????? ?????????????????????????? ?????????????? ???????????????????? ???? ??????????????????????")
	    }
	    return nil
    }
}

// Validate validates the QuickTodoForm fields
func (m QuickSearchForm) Validate() error {
	return vld.ValidateStruct(&m,
		vld.Field(&m.Clause, vld.Required, vld.Length(3, 30), vld.By(swearDetector()), vld.Match(regexp.MustCompile(`^(([0-9a-zA-Z' -]{3,30})|([0-9??-????-??????????????????????' -]{3,30}))`))),
	)
}

// Validate validates the WatchAuthorForm fields
func (m WatchAuthorForm) Validate() error {
	return vld.ValidateStruct(&m,
		vld.Field(&m.TodoId, vld.Required),
		//!!! vld.Field(&m.RecaptchaResponse, vld.When(config.CFG.AppMode != "dev", vld.Required, vld.Length(64, 512)).Else(vld.NotNil)),
		vld.Field(&m.RecaptchaResponse, vld.When(false, vld.Required, vld.Length(64, 512)).Else(vld.NotNil)),
		vld.Field(&m.SignerUA, vld.Required, vld.Length(30,512), vld.Match(regexp.MustCompile(`[0-9a-zA-Z-\/;\(\)\.,]{30,512}`))),
		vld.Field(&m.SignerScreen, vld.Required, vld.Length(2, 64), vld.Match(regexp.MustCompile(`[0-9]{3,5}x[0-9]{3,5}`))),
		vld.Field(&m.SignerLangs, vld.Required, vld.Length(5, 256), vld.Match(regexp.MustCompile(`[a-zA-Z,-]{5,256}`))),
		vld.Field(&m.SignerTime, vld.Required, vld.Length(12, 16), vld.Match(regexp.MustCompile(`[\d]{12,16}`))),
	)
}

// Validate validates the QuickTodoForm fields
func (m LoginForm) Validate() error {
	return vld.ValidateStruct(&m,
		vld.Field(&m.Username, vld.When(config.CFG.AppMode != "dev", vld.Required, is.Email).Else(vld.Required, is.EmailFormat)),
		vld.Field(&m.CurrentPassword, vld.Required, vld.Length(6, 128)),
	)
}

// Validate validates the ProfileForm fields
func (m ProfileForm) Validate() error {
	return vld.ValidateStruct(&m,
		vld.Field(&m.Tel, vld.When(m.Tel != "", vld.Required, vld.Length(10, 21), vld.Match(regexp.MustCompile(`[\d\s\-\+\(\)]{10,21}`))).Else(vld.Nil)),
		vld.Field(&m.GivenName, vld.Required, vld.Length(2, 64), vld.Match(regexp.MustCompile(`^(([a-zA-Z' -]{2,128})|([??-????-??????????????????????' -]{2,128}))`))),
		vld.Field(&m.NewPassword, vld.When(m.NewPassword != "", vld.Length(6, 128), vld.Required)),
		vld.Field(&m.NewPasswordRepeat, vld.When(m.NewPassword != "", vld.By(passwordsEquals(m.NewPassword)))),
	)
}

func (m MessageForm) Validate() error {
	return vld.ValidateStruct(&m,
		vld.Field(&m.Email, vld.When(config.CFG.AppMode != "dev", vld.Required, is.Email).Else(vld.Required, is.EmailFormat)),
		vld.Field(&m.GivenName, vld.Required, vld.Length(2, 64), vld.Match(regexp.MustCompile(`^(([a-zA-Z' -]{2,128})|([??-????-??????????????????????' -]{2,128}))`))),
		vld.Field(&m.Msg, vld.Length(12, 512)),
	)
}

func (m SupportForm) Validate() error {
	return vld.ValidateStruct(&m,
		vld.Field(&m.Email, vld.When(config.CFG.AppMode != "dev", vld.Required, is.Email).Else(vld.Required, is.EmailFormat)),
		vld.Field(&m.GivenName, vld.Required, vld.Length(2, 64), vld.Match(regexp.MustCompile(`^(([a-zA-Z' -]{2,128})|([??-????-??????????????????????' -]{2,128}))`))),
		vld.Field(&m.Subject, vld.Length(3, 12)),
		vld.Field(&m.Msg, vld.Length(12, 512)),
	)
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
	return ag.repo.UpdateUser(ctx, entity.User{
		// UserId: uid,
    Name:     pf.GivenName,
    // Email:    pf.Email,
    Tel:      util.PhoneNormalisation(pf.Tel),
    // Impp:     "",
    // Authkey:  util.GetSaltedSha256(config.CFG.AppSecretKey, pf.Email),
    Passhash: pf.NewPassword,
    // Approvetoken: "",
    Picture:   avafile,
    // Created:   time.Now().Format("2006-01-02 15:04:05"),
    Lastlogin: time.Now().Format("2006-01-02 15:04:05"),
    // Roles:     "user",
    // Notes:     "regform",
	}, uid)
}

func (ag agregator) UpdateTodoPicture(ctx context.Context, aid int64, field, fname string) error {
	return ag.repo.UpdateTodoPicture(ctx, aid, field, fname)
}

func (ag agregator) UpdateUserLastlogin(ctx context.Context, uid int64, dtstring string) error {
	return ag.repo.UpdateUserLastlogin(ctx, uid, dtstring)
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
	Todo, err := ag.repo.GetTodoById(ctx, id)
	if err != nil {
		return Todo{}, err
	}
	return Todo{Todo}, nil
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

func (ag agregator) GetMessagesSendersByUserId(ctx context.Context, uid int64) ([]MessageSender, error) {
	items, err := ag.repo.GetMessagesSendersByUserId(ctx, uid)
	if err != nil {
		return nil, err
	}
	result := []MessageSender{}
	for _, item := range items {
		result = append(result, MessageSender{item})
	}
	return result, nil
}

func (ag agregator) GetTodoLast(ctx context.Context) ([]Todo, error) {
	items, err := ag.repo.GetTodoLast(ctx)
	if err != nil {
		return nil, err
	}
	result := []Todo{}
	for _, item := range items {
		result = append(result, Todo{item})
	}
	return result, nil
}

func (ag agregator) GetTodoSearch(ctx context.Context, clause string) ([]Todo, error) {
	items, err := ag.repo.GetTodoSearch(ctx, clause)
	if err != nil {
		return nil, err
	}
	result := []Todo{}
	for _, item := range items {
		result = append(result, Todo{item})
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
	return ag.repo.DeleteTodoById(ctx, aid)
}

// Create creates a new Todo.
func (ag agregator) CreateTodo(ctx context.Context, f *QuickTodoForm, uid int64, dt string) (int64, error) {
	if uid == 0 {
		return 0, errors.New("CreateTodo User_id:0")
	}
	return ag.repo.CreateTodo(ctx, entity.Todo{
	  UserId: uid,
	  CategoryId: util.Pkeyer(f.CategoryId),
	  Title: util.ExtractTitle(f.Nanopost, 45, 4),
	  Nanopost: f.Nanopost,
	  Price: util.ExtractDigitsInt(f.Price),
	  Currency: config.CURRENCY,
	  Picture1: "",
	  Picture2: "",
	  Picture3: "",
	  Picture4: "",
	  // Picture5: "",
	  // Picture6: "",
	  // ModeratorId: 0,
	  Created: dt,
	  Active: 1,
	})
}

// Create creates a new album.
func (ag agregator) CreateMessage(ctx context.Context, fromId, toId int64, msg, dtstring string) (error) {
	return ag.repo.CreateMessage(ctx, entity.Message{
	  SenderId:   fromId,
	  ReceiverId: toId,
	  Content:    msg,
	  Sended:     dtstring, // Datetime for getting support picture files
	  Readed:     "0000-00-00 00:00:00",
		ModeratorId:0,
	})
}

// func (ag agregator) GetCategories(ctx context.Context) ([]Category, error) {
func (ag agregator) GetCategoriesPath(ctx context.Context) ([]CategoryPath, error) {
	items, err := ag.repo.GetCategoryPath(ctx)
	if err != nil {
		return nil, err
	}
	result := []CategoryPath{}
	for _, item := range items {
		result = append(result, CategoryPath{item})
	}
	return result, nil
}
