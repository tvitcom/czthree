package todos

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
