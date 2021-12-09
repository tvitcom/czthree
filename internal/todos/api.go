package Todo

import (
    "time"
	"github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/limiter"
    "github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/tvitcom/czthree/pkg/log"
    "github.com/tvitcom/czthree/pkg/util"
    "github.com/tvitcom/czthree/internal/config"
    "github.com/valyala/fastjson"
)

type (
    resource struct {
    	agregator Agregator
    	logger  log.Logger
    }
) 

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(router *fiber.App, agregator Agregator, logger log.Logger) {
	res := resource{agregator, logger}
    MWheader := func(c *fiber.Ctx) error {
        c.Append("Powered-by", config.CFG.WebservName)
        
        // if config.CFG.AppMode != "dev" {
        //    c.Append("Strict-Transport-Security", "max-age=7776000; includeSubDomains")
        //    c.Set("X-XSS-Protection", "1; mode=block")
        //    c.Set("X-Content-Type-Options", "nosniff")
        //    c.Set("X-Download-Options", "noopen")
        //    c.Set("Strict-Transport-Security", "max-age=5184000")
        //    c.Set("X-Frame-Options", "SAMEORIGIN")
        //    c.Set("X-DNS-Prefetch-Control", "off")
        // }
          // Set some security headers:
        return c.Next()
    }
    
    MWauthentication := func(c *fiber.Ctx) error {
        c.Append("Restricted-by", "jwt")
        tryjwt := c.Cookies("tok","~")
        // Cookie empty
        if tryjwt == "" {
            return c.Status(403).Redirect("/error.html?msg=Используйте страницу входа")
        }
        // Cookie with any jwt string
        tokdata, _, err := util.GetJwtClaimHMAC(tryjwt, config.CFG.AppSecretKey)
        if err != nil {
            println("AUTH:err:", err.Error())
            return c.Status(403).Redirect("/error.html?msg=Ошибка авторизации")

        }
        _ = c.Locals("iam", fastjson.GetString(tokdata, "sub"))
        return c.Next()
    }
    
    MWuserinfo := func(c *fiber.Ctx) error {
        tryjwt := c.Cookies("tok","")
        // Cookie empty
        if tryjwt != "" {
            tokdata, _, errtok := util.GetJwtClaimHMAC(tryjwt, config.CFG.AppSecretKey)
            if errtok != nil {            
                return c.Status(403).Redirect("/error.html?msg=Ошибка маркера доступа")
            }
            _ = c.Locals("iam", fastjson.GetString(tokdata, "sub"))
        }
        return c.Next()
    }
    
    MWnouser := func(c *fiber.Ctx) error {
        mwuid := c.Locals("iam")
        if mwuid != nil {
            return c.Redirect("/my/todolist.html", 301)
        }
        return c.Next()
    }
    
    MWcsp := func(c *fiber.Ctx) error {
            // require-trusted-types-for 'script';
        csp := `
            default-src 'self';
            connect-src 'self' https://www.google-analytics.com https://www.google.com/recaptcha/ https://www.gstatic.com/recaptcha/;
            font-src 'self' https://fonts.gstatic.com;
            frame-src 'self' https://www.google.com/recaptcha/ https://www.google.com/maps/ https://youtu.be https://youtube.com https://www.youtube.com;
            frame-ancestors https://youtu.be https://youtube.com https://www.youtube.com;
            img-src 'self' https://www.google.com/recaptcha/ https://lh3.googleusercontent.com/ https://images.unsplash.com data: blob: https://source.unsplash.com;
            object-src 'none';
            script-src 'self' 'unsafe-inline' 'unsafe-eval' https://www.google.com https://apis.google.com https://www.gstatic.com/recaptcha/ https://www.googletagmanager.com https://www.google-analytics.com;
            style-src 'self' 'unsafe-inline' https://fonts.googleapis.com;
            report-uri https://` + config.CFG.AppFqdn + `/csp_collector.html
        `
            // if c.Request.Method == "OPTIONS" {
            //     if len(c.Request.Header["Access-Control-Request-Headers"]) > 0 {
            //         c.Header("Access-Control-Allow-Headers", c.Request.Header["Access-Control-Request-Headers"][0])
            //     }
            //     c.AbortWithStatus(http.StatusOK)
            // }
        var policy string
        if config.CFG.AppMode == "dev" {
            policy = `default-src 'self' 'unsafe-inline' 'unsafe-eval' data: blob:;
                img-src 'self' data: blob:;
                object-src 'self';
                script-src 'self' 'unsafe-inline' 'unsafe-eval';
                style-src 'self' 'unsafe-inline';`
        } else {
            policy = csp
        }
        c.Append("Content-Security-Policy", policy)
        c.Append("X-Content-Type-Options", "nosniff")
        c.Append("X-Frame-Options", "SAMEORIGIN")
        return c.Next()
    }
    
    MWRateLim := limiter.New(limiter.Config{
        Next: func(c *fiber.Ctx) bool {
            return c.IP() == "127.0.0.1"
        },
        Max:          1,
        Expiration:   1 * time.Minute,
        KeyGenerator: func(c *fiber.Ctx) string {
            return c.Get("x-forwarded-for")
        },
        LimitReached: func(c *fiber.Ctx) error {
            c.Status(fiber.StatusTooManyRequests)
            return c.Render("error", fiber.Map{
                "msg": "Быстрее только кролики...",
            })
        },
    })
	
    
// GET  /index.html?loc=Kharkovskaya&cat=123&q=qwerty
    // router.Use(MWRateLim)
    router.Use(compress.New(compress.Config{
        Level: compress.LevelBestSpeed,
    }))
    
	router.Get("/", MWuserinfo, MWcsp, res.pageIndex)
    // GET  /error.html?m=You%20will%20signup%20firstly
    router.Get("/error.html", MWuserinfo, MWheader, MWcsp, res.pageError)
    
// RATELIMITED PAGES:
    authGroup := router.Group("/auth", MWRateLim, MWuserinfo, MWcsp)

// GET  /auth/login.html?uid=1234&aprove=mail&otcode=123qwerty5467
	authGroup.Get("/login.html", MWnouser, res.pageLogin)

// POST /auth/login.html?aprove=mail
	authGroup.Post("/login.html", MWnouser, res.handlerLogin)

// POST /auth/logout.html
    authGroup.Post("/logout.html", res.handlerLogout)

// REGISTERED USERS ONLY:
//Создание/Просмотр/Изменение/Удаление TODO.
    myGroup := router.Group("/my", MWauthentication, MWheader, MWcsp)

/* ----------- Todo: --------------*/
// GET  /my/todolist.html&uid=123
    myGroup.Get("/todolist.html", res.pageUserTodo)
    myGroup.Post("/todochange.html", res.handlerUpdateTodo)

// POST /my/todoupdate.html&tid=123
    // myGroup.Post("/todoupdate.html", res.handlerUpdateTodo)
// POST /my/tododelete.html&tid=123
    // myGroup.Post("/tododelete.html", res.handlerDeleteTodo)

// GET /my/todonew.html
    // myGroup.Get("/todonew.html", res.pageSoon)
// POST /my/todonew.html
    // myGroup.Post("/todonew.html", res.pageSoon)

/* ---------USERS: ---------------- */
    myGroup.Get("/userlist.html", res.pageUserList)

    // POST /my/userdelete.html?uid=123
    myGroup.Post("/userlist.html", res.pageSoon)

// GET  /my/userprofile.html?uid=123
    myGroup.Get("/userprofile.html", res.pageUserProfile)
// POST /my/userprofile.html
    myGroup.Post("/userprofile.html", res.handlerUserProfile)
}