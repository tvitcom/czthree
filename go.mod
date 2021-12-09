module github.com/tvitcom/czthree

go 1.16

replace github.com/tvitcom/czthree/internal/errors => ./internal/errors

replace github.com/tvitcom/czthree/pkg/log => ./pkg/log

replace github.com/tvitcom/czthree/pkg/dbcontext => ./pkg/dbcontext

replace github.com/tvitcom/czthree/internal/config => ./internal/config

replace github.com/tvitcom/czthree/pkg/pagination => ./pkg/pagination

replace github.com/tvitcom/czthree/internal/entity => ./internal/entity

replace github.com/tvitcom/czthree/internal/dto => ./internal/dto

replace github.com/tvitcom/czthree/internal/healthcheck => ./internal/healthcheck

replace github.com/tvitcom/czthree/internal/user => ./internal/user

replace github.com/tvitcom/czthree/internal/todos => ./internal/todos

replace github.com/tvitcom/czthree/internal/test => ./internal/test

replace github.com/tvitcom/czthree/pkg/accesslog => ./pkg/accesslog

replace github.com/tvitcom/czthree/bkp/swearfilter => ./pkg/swearfilter

require (
	github.com/alexedwards/argon2id v0.0.0-20211130144151-3585854a6387
	github.com/dvsekhvalnov/jose2go v1.5.0
	github.com/go-ozzo/ozzo-dbx v1.5.0
	github.com/go-ozzo/ozzo-routing/v2 v2.4.0 // indirect
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0 // indirect
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gofiber/fiber/v2 v2.23.0
	github.com/gofiber/template v1.6.20
	github.com/nats-io/nats-server/v2 v2.6.6 // indirect
	github.com/nats-io/nats.go v1.13.1-0.20211122170419-d7c1d78a50fc
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/qiangxue/go-env v1.0.1 // indirect
	github.com/rwcarlsen/goexif v0.0.0-20190401172101-9e8deecbddbd
	github.com/tvitcom/czthree/bkp/swearfilter v0.0.0-00010101000000-000000000000
	github.com/tvitcom/czthree/internal/config v0.0.0-00010101000000-000000000000
	github.com/tvitcom/czthree/internal/dto v0.0.0-00010101000000-000000000000 // indirect
	github.com/tvitcom/czthree/internal/entity v0.0.0-00010101000000-000000000000 // indirect
	github.com/tvitcom/czthree/internal/errors v0.0.0-00010101000000-000000000000 // indirect
	github.com/tvitcom/czthree/internal/healthcheck v0.0.0-00010101000000-000000000000
	github.com/tvitcom/czthree/internal/test v0.0.0-00010101000000-000000000000 // indirect
	github.com/tvitcom/czthree/internal/todos v0.0.0-00010101000000-000000000000
	github.com/tvitcom/czthree/pkg/accesslog v0.0.0-00010101000000-000000000000 // indirect
	github.com/tvitcom/czthree/pkg/dbcontext v0.0.0-00010101000000-000000000000
	github.com/tvitcom/czthree/pkg/log v0.0.0-00010101000000-000000000000
	github.com/valyala/fastjson v1.6.3 // indirect
	golang.org/x/crypto v0.0.0-20211202192323-5770296d904e
)
