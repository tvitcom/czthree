package todos

import (
    // "fmt"
    // "time"
    // "errors"
    "github.com/gofiber/fiber/v2"
    "github.com/tvitcom/czthree/internal/config"
    "github.com/tvitcom/czthree/pkg/util"  
    
)

func (res resource) pageUserProfile(c *fiber.Ctx) error {
    uid := util.Pkeyer(c.Locals("iam"))
    user, err := res.agregator.GetUserById(c.UserContext(), uid)
    if err != nil {
        return err
    }
    return c.Render("userprofile", fiber.Map{
        "msg": "редактирование данных",
        "user": user,
    })
}

func (res resource) handlerUserProfile(c *fiber.Ctx) error {
    uid := util.Pkeyer(c.Locals("iam"))
    curruser, err := res.agregator.GetUserById(c.UserContext(), uid)
    if err != nil {
        return c.Status(412).Render("error", fiber.Map{"msg": err.Error()})
    }

    form := new(ProfileForm)
    if err := c.BodyParser(form); err != nil {
        res.logger.With(c.UserContext()).Error(err.Error())
        return c.Status(412).Render("error", fiber.Map{
            "msg": err.Error(),
            "user": curruser,
        })
    }

    // form.Tel = util.PhoneNormalisation(form.Tel)   
    // Form validation    
    if err := form.Validate(); err != nil {
        return c.Status(412).Render("error", fiber.Map{"msg": err.Error()})
    }

    err = res.agregator.UpdateUser(c.UserContext(), form, uid)
    if err != nil {
        return c.Status(412).Render("error", fiber.Map{"msg": err.Error()})
    }
    return c.Redirect("/my/userprofile.html")
}

func (res resource) pageUserTodo(c *fiber.Ctx) error {
    uid := util.Pkeyer(c.Locals("iam"))
    var perfomerId int64
    perfomerId = util.Pkeyer(c.Query("uid", "0"))
    if uid != config.AdminUserID && perfomerId != 0 {
        return c.Status(403).Render("error", fiber.Map{"msg": "Unauthorised"})
    }
    curruser, err := res.agregator.GetUserById(c.UserContext(), uid)
    if err != nil {
        return c.Status(412).Render("error", fiber.Map{"msg": err.Error()})
    }
    if perfomerId == 0 {
        perfomerId = uid
    }
    perfomer, err := res.agregator.GetUserById(c.UserContext(), perfomerId)
    if err != nil {
        return c.Status(412).Render("error", fiber.Map{"msg": err.Error()})
    }
    todosDisplay, err := res.agregator.GetTodoDisplayByUserId(c.UserContext(), perfomerId)
    if err != nil {
        println("DBERR:", err.Error())
        return c.Status(500).Redirect("/error.html?msg=Ошибка работы сайта")
    }
    return c.Render("usertodolist", fiber.Map{
        "msg": "userTodo page: page: Coming soon!",
        "todosDisplay": todosDisplay,
        "perfomer": perfomer,
        "user": curruser,
    })
}

func (res resource) handlerUpdateTodoStatus(c *fiber.Ctx) error {
    uid := util.Pkeyer(c.Locals("iam"))
    form := new(TodoStatusForm)
    form.TodoId = util.Pkeyer(c.Query("todoid", "0"))
    if err := c.BodyParser(form); err != nil && uid != 0 {
        return c.Status(500).Redirect("/error.html?msg=Ошибка обработки формы для смены статуса todo")
    }
    // Form validation
    if err := form.Validate(); err != nil {
        return c.Status(412).Render("error", fiber.Map{"msg": err.Error()})
    }
    if err := res.agregator.UpdateTodoStatus(c.UserContext(), form, uid); err != nil {
        res.logger.With(c.UserContext()).Error(err.Error())
        return c.Status(500).Redirect("/error.html?msg=Ошибка редактирования todo")
    }
    return c.Status(302).Redirect("/my/todolist.html")
}

func (res resource) pageUpdateTodo(c *fiber.Ctx) error {
    uid := util.Pkeyer(c.Locals("iam"))
    if uid != config.AdminUserID {
        return c.Status(403).Render("error", fiber.Map{"msg": "Unauthorised"})
    }
    // -form := new(TodoForm)
    todoId := util.Pkeyer(c.Query("todoid", "0"))
    todo, err := res.agregator.GetTodoById(c.UserContext(), todoId)
    if err != nil {
        return c.Status(412).Render("error", fiber.Map{"msg": err.Error()})
    }
    curruser, err := res.agregator.GetUserById(c.UserContext(), uid)
    if err != nil {
        return c.Status(412).Render("error", fiber.Map{"msg": err.Error()})
    }
    users, err := res.agregator.GetUsersWithLimitOffset(c.UserContext(), 100, 0)
    if err != nil {
        return c.Status(412).Render("error", fiber.Map{"msg": err.Error()})
    }
    return c.Status(200).Render("usertodoform", fiber.Map{
        "todo": todo,
        "user": curruser,
        "users": users,
    })
}

func (res resource) handlerUpdateTodo(c *fiber.Ctx) error {
    uid := util.Pkeyer(c.Locals("iam"))
    if uid != config.AdminUserID {
        return c.Status(403).Render("error", fiber.Map{"msg": "Unauthorised"})
    }
    var perfomerId int64
    perfomerId = util.Pkeyer(c.Query("uid", "0"))
    if perfomerId == 0 {
        perfomerId = uid
    }
    form := new(TodoForm)    
    if err := c.BodyParser(form); err != nil && uid != 0 {
        return c.Status(500).Redirect("/error.html?msg=Ошибка обработки формы для смены статуса todo")
    }
    // Form validation
    if err := form.Validate(); err != nil {
        return c.Status(412).Render("error", fiber.Map{"msg": err.Error()})
    }
    if err := res.agregator.UpdateTodo(c.UserContext(), form); err != nil {
        res.logger.With(c.UserContext()).Error(err.Error())
        return c.Status(500).Redirect("/error.html?msg=Ошибка редактирования todo")
    }
    return c.Status(302).Redirect("/my/userlist.html?uid")
}

func (res resource) pageCreateTodo(c *fiber.Ctx) error {
    uid := util.Pkeyer(c.Locals("iam"))
    if uid != config.AdminUserID {
        return c.Status(403).Render("error", fiber.Map{"msg": "Unauthorised"})
    }
    curruser, err := res.agregator.GetUserById(c.UserContext(), uid)
    if err != nil {
        return c.Status(412).Render("error", fiber.Map{"msg": err.Error()})
    }
    users, err := res.agregator.GetUsersWithLimitOffset(c.UserContext(), 100, 0)
    if err != nil {
        return c.Status(412).Render("error", fiber.Map{"msg": err.Error()})
    }
    return c.Status(200).Render("usertodonew", fiber.Map{
        "user": curruser,
        "users": users,
    })
}

func (res resource) handlerCreateTodo(c *fiber.Ctx) error {
    uid := util.Pkeyer(c.Locals("iam"))
    if uid != config.AdminUserID {
        return c.Status(403).Render("error", fiber.Map{"msg": "Unauthorised"})
    }
    form := new(TodoForm)    
    if err := c.BodyParser(form); err != nil {
        return c.Status(500).Redirect("/error.html?msg=Ошибка обработки формы для создания todo")
    }
    // Form validation
    if err := form.Validate(); err != nil {
        return c.Status(412).Render("error", fiber.Map{"msg": err.Error()})
    }
    if i, err := res.agregator.CreateTodo(c.UserContext(), form); err != nil || i == 0 {
        res.logger.With(c.UserContext()).Error(err.Error())
        return c.Status(500).Redirect("/error.html?msg=Ошибка создания todo")
    }
    return c.Status(302).Redirect("/my/todolist.html")
}

func (res resource) handlerDeleteTodo(c *fiber.Ctx) error {
    uid := util.Pkeyer(c.Locals("iam"))
    if uid != config.AdminUserID { // Non admin with id=1 go ahead
        return c.Status(403).Render("error", fiber.Map{"msg": "Unauthorised"})
    }
    form := new(DeleteTodoForm)
    if err := c.BodyParser(form); err != nil && uid != 0 {
        return c.Status(500).Redirect("/error.html?msg=Ошибка обработки формы для удаления объявления")
    }
    if err := res.agregator.DeleteTodoData(c.UserContext(), form.TodoId); err != nil {
        res.logger.With(c.UserContext()).Error(err.Error())
        return c.Status(500).Redirect("/error.html?msg=Ошибка удаления записи запланированного дела")
    }
    return c.Status(201).Render("thanks", fiber.Map{
        "msg": "Дело успешно удалено",
    })
}

func (res resource) pageUserList(c *fiber.Ctx) error {
    uid := util.Pkeyer(c.Locals("iam"))
    if uid != config.AdminUserID { // Non admin with id=1 go ahead
        return c.Status(403).Render("error", fiber.Map{"msg": "Unauthorised"})
    }
    user, err := res.agregator.GetUserById(c.UserContext(), uid)
    if err != nil {
        return c.Status(412).Render("error", fiber.Map{"msg": err.Error()})
    }
    users, err := res.agregator.GetUsersWithLimitOffset(c.UserContext(), 1000, 0)
    if err != nil {
        return c.Status(412).Render("error", fiber.Map{"msg": err.Error()})
    }
    return c.Render("userlist", fiber.Map{
        "msg": "userlist page: page: Coming soon!",
        "user": user,
        "users": users,
    })
}
