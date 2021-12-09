package Todo

import (
    // "strings"
    // "errors"
    // "time"
    "github.com/gofiber/fiber/v2"
    "github.com/tvitcom/czthree/internal/config"
    // "github.com/tvitcom/czthree/pkg/util"  
)

func (res resource) pageIndex(c *fiber.Ctx) error {
//     var Todo []Todo
//     searchclause := strings.TrimSpace(c.Query("q", ""))
    return c.Render("index", fiber.Map{
//         "msg": "index page",
//         "Todo": Todo,
    })
}

func (res resource) pageWatch(c *fiber.Ctx) error {
    // uid := util.Pkeyer(c.Locals("iam"))
    // todoid := util.Pkeyer(c.Query("todoid", "0"))
    // if todoid == 0 {
    //     res.logger.With(c.UserContext()).Error("Ошибка GET todoid параметра")
    //     return c.Status(500).Redirect("/error.html?msg=Ошибка параметра запроса.")
    // }
    // todo, err := res.agregator.GetTodoById(c.UserContext(), todoid)
    // if err != nil {
    // println(err.Error())
    //     return c.Status(404).Redirect("/error.html?msg=Объявление не найдено или удалено.")
    // }
    // author, err := res.agregator.GetUserById(c.UserContext(), todo.AuthorId)
    // if err != nil {
    //     return c.Status(500).Redirect("/error.html?msg=Ошибка работы сайта.")
    // }
    return c.Render("watch", fiber.Map{
        "msg": "watch page",
    //     "recaptcha_site_key": config.CFG.RecaptchaSiteKey,
    //     "todo": todo,
    //     "user": author,
    //     "fqdn": config.CFG.AppFqdn,
    //     "uid": uid,
    })
}

func (res resource) handlerActivity(c *fiber.Ctx) error {
    // uid := util.Pkeyer(c.Query("uid", "0"))
    // if uid == 0 {
    //     res.logger.With(c.UserContext()).Error("Ошибка GET todoid параметра")
    //     return c.Status(500).Redirect("/error.html?msg=Ошибка параметра")
    // }
    // Todo, err := res.agregator.GetTodoByUserId(c.UserContext(), uid)
    // if err != nil {
    //     return c.Status(500).Redirect("/error.html?msg=Ошибка работы сайта")
    // }
    // author, err := res.agregator.GetUserById(c.UserContext(), uid)
    // if err != nil {
    //     return c.Status(500).Redirect("/error.html?msg=Ошибка работы сайта")
    // }
    return c.Render("activity", fiber.Map{
        // "msg": "activity page",
        // "Todo": Todo,
        // "author": author,
    })
}

// ->makeUser->makeTodoInactive->sendApproveUser1 [->handleEmailApprove]->userprofile->userTodo


func (res resource) pageError(c *fiber.Ctx) error {  
    return c.Render("error", fiber.Map{
        "msg": c.Query("msg", "Чтото пошло не так."),
    })
}

func (res resource) pageThanks(c *fiber.Ctx) error {
    return c.Render("thanks", fiber.Map{
            "thanks": c.Query("msg", "Благодарим за использование нашего сайта " + config.CFG.BizName),
        })
}

func (res resource) pageSoon(c *fiber.Ctx) error {
    return c.Render("soon", fiber.Map{
        "msg": "Контент станет доступен позже.",
    })
}

func (res resource) pageAgreement(c *fiber.Ctx) error {
    return c.SendFile("./assets/docs/agreement.txt", true)
}

func (res resource) pageGdprPolicy(c *fiber.Ctx) error {
    return c.SendFile("./assets/docs/GDPR_POLICY_RU.txt", true)
}

func (res resource) pageRobots(c *fiber.Ctx) error {
    return c.SendFile(config.RobotsFilePath, true)
}

func (res resource) pageSitemap(c *fiber.Ctx) error {
    return c.SendFile(config.SitemapFilePath, true)
}

func (res resource) handlerCspCollector(c *fiber.Ctx) error {
    res.logger.With(c.UserContext()).Info(string(c.Body()))
    return c.JSON(&fiber.Map{
            "ok": true,
            "data": "log ok",
        },
    )
}
