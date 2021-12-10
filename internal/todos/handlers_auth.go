package todos

import (
    "time"
    // "fmt"
    // "errors"
    "github.com/gofiber/fiber/v2"
    // "github.com/tvitcom/czthree/pkg/log"
    "github.com/tvitcom/czthree/internal/config"
    "github.com/tvitcom/czthree/pkg/util"  
)


func (res resource) pageLogin(c *fiber.Ctx) error {
    return c.Render("login", fiber.Map{
        "msg": "avdert page: Coming soon",
    })
}

func (res resource) handlerLogin(c *fiber.Ctx) error {
    // handle_dt := time.Now().Format("2006-01-02 15:04:05")
    form := new(LoginForm)
    if err := c.BodyParser(form); err != nil {
        res.logger.With(c.UserContext()).Error(err.Error())
        return c.Status(412).Render("error", fiber.Map{"msg": err.Error()})
    }
    // validation    
    if err := form.Validate(); err != nil {
        return c.Status(403).Render("error", fiber.Map{
            "msg": "Неверно введены данные формы логина",
        })
    }
    
    user, err := res.agregator.GetUserByEmail(c.UserContext(), form.Username)
 
    // currentHash := util.MakeBCryptHash(form.CurrentPassword, config.BCRYPT_COST)
    // fmt.Println("THAT PASSWORD:", currentHash)
    // fmt.Printf("SQL: UPDATE user SET passhash = '%s' WHERE email = '%s';\n", currentHash, form.Username)
    
    if err != nil || util.VerifyBCrypt(form.CurrentPassword, user.Passhash) != nil {
        return c.Status(403).Render("error", fiber.Map{
            "msg": "Неверно ввели логин или пароль или всё вместе",
        })
    }
    // Update the user.lastlogin
    // if err := res.agregator.UpdateUserLastlogin(c.UserContext(), user.UserId, handle_dt); err != nil {
    //     return c.Status(500).Redirect("/error.html?msg=Ошибка работы сайта")
    // }
    
    // Let identity marker - user authenticated successfully
    // seckey, tokid, fqdn, uid, appsid, roles
    rnd32 := util.RandomHexString(8)
    tok, err := util.MakeJwtString(config.CFG.AppSecretKey, rnd32, config.CFG.AppFqdn, util.Stringer(user.UserId), "main", "user")
    if err != nil {
        return c.Status(403).Render("error", fiber.Map{"msg": err.Error()})
    }
    makeJWTCookie(c, tok)
    return c.Redirect("/my/todolist.html", 301)
}
func (res resource) handlerLogout(c *fiber.Ctx) error {
    deleteJWTCookie(c)
    return c.Render("thanks", fiber.Map{
        "msg": "за посещение сайта. Удачных сделок!",
    })
}

func makeJWTCookie(c *fiber.Ctx, jwt string) {
    duration := 120 * time.Minute
    cookie := new(fiber.Cookie)
    cookie.Name = "tok"
    cookie.Value = jwt
    cookie.HTTPOnly = true
    cookie.Expires = time.Now().Add(duration)
    c.Cookie(cookie)
}

func deleteJWTCookie(c *fiber.Ctx) {
    negateDuration := -1 * time.Minute
    cookie := new(fiber.Cookie)
    cookie.Name = "tok"
    cookie.Value = "logout"
    cookie.HTTPOnly = true
    cookie.Expires = time.Now().Add(negateDuration)
    c.Cookie(cookie)
}