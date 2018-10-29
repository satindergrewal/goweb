package main

import (
	"fmt"
	"github.com/kataras/iris"
	"strconv"

	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

type Add struct {
	Num1 int `json:"num1"`
	Num2 int `json:"num2"`
}

func add(val1 int, val2 int) int {
	result_val := val1 + val2
	//fmt.Println(result_val)
	return result_val
}

type User struct {
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	City      string `json:"city"`
	Age       int    `json:"age"`
}

func main() {
	res := add(3, 5)
	fmt.Println(res)

	app := iris.New()
	app.Logger().SetLevel("debug")
	// Optionally, add two built'n handlers
	// that can recover from any http-relative panics
	// and log the requests to the terminal.
	app.Use(recover.New())
	app.Use(logger.New())

	// Method:   GET
	// Resource: http://localhost:8080
	app.Handle("GET", "/", func(ctx iris.Context) {
		ctx.HTML("<h1>Welcome</h1>")
	})

	app.Handle("GET", "/contact", func(ctx iris.Context) {
		ctx.HTML("<h1> Hello from /contact </h1>")
	})

	// same as app.Handle("GET", "/ping", [...])
	// Method:   GET
	// Resource: http://localhost:8080/ping
	app.Get("/ping", func(ctx iris.Context) {
		ctx.WriteString("pong")
	})

	// Method:   GET
	// Resource: http://localhost:8080/hello
	app.Get("/hello", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"message": "Hello Iris!"})
	})

	app.Favicon("./assets/favicon.ico")

	// enable gzip, optionally:
	// if used before the `StaticXXX` handlers then
	// the content byte range feature is gone.
	// recommend: turn off for large files especially
	// when server has low memory,
	// turn on for medium-sized files
	// or for large-sized files if they are zipped already,
	// i.e "zippedDir/file.gz"
	//
	app.Use(iris.Gzip)

	// first parameter is the request path
	// second is the system directory
	//
	app.StaticWeb("/css", "./assets/css")
	app.StaticWeb("/js", "./assets/js")
	//
	app.StaticWeb("/static", "./assets")

	// http://localhost:8080/static/css/main.css
	// http://localhost:8080/static/js/jquery-2.1.1.js
	// http://localhost:8080/static/favicon.ico

	// http://localhost:8080
	// http://localhost:8080/ping
	// http://localhost:8080/hello

	app.Get("/profile/{username:string}", info)
	app.Get("/profile/{username}/backups/{filepath:path}", info)

	// GET: http://localhost:8080/profile/anyusername
	// GET: http://localhost:8080/profile/anyusername/backups/any/number/of/paths/here

	usersRoutes := app.Party("/users")
	// GET: http://localhost:8080/users/help
	usersRoutes.Get("/help", func(ctx iris.Context) {
		ctx.Writef("POST / -- create new user\n")
	})

	// POST: http://localhost:8080/users
	usersRoutes.Post("/", func(ctx iris.Context) {
		username, password := ctx.PostValue("username"), ctx.PostValue("password")
		ctx.Writef("create user for username= %s and password= %s", username, password)
	})

	// Method POST: http://localhost:8080/decode
	app.Post("/decode", func(ctx iris.Context) {
		var user User
		ctx.ReadJSON(&user)
		ctx.Writef("%s %s is %d years old and comes from %s", user.Firstname, user.Lastname, user.Age, user.City)
	})

	// Method GET: http://localhost:8080/encode
	app.Get("/encode", func(ctx iris.Context) {
		doe := User{
			Username:  "Johndoe",
			Firstname: "John",
			Lastname:  "Doe",
			City:      "Neither FBI knows!!!",
			Age:       25,
		}

		ctx.JSON(doe)
	})


	mathsRoutes := app.Party("/maths")
	// GET: http://localhost:8080/maths/help
	mathsRoutes.Get("/help", func(ctx iris.Context) {
		ctx.Writef("POST / -- mathematics cool as!\n")
	})

	// GET: http://localhost:8080/maths/add
	mathsRoutes.Post("/add", func(ctx iris.Context) {
		//num1, _ := ctx.Params().GetInt("num1")
		//num1 := ctx.PostValue("num1")
		num1, err := strconv.Atoi(ctx.PostValue("num1"))
		if err != nil {
			fmt.Println(err)
		}
		//num2 := ctx.PostValue("num2")
		num2, err := strconv.Atoi(ctx.PostValue("num2"))
		if err != nil {
			fmt.Println(err)
		}
		math_res := add(num1, num2)
		ctx.Writef("Addition of %d + %d is: %d", num1, num2, math_res)
		//ctx.Writef("got number: %T and %T", num1, num2)
	})

	calcRoutes := app.Party("/calc")
	// GET: http://localhost:8080/calc/help
	calcRoutes.Get("/help", func(ctx iris.Context) {
		ctx.Writef("POST / -- calculator helper\n")
	})

	app.Run(iris.Addr(":8080"), iris.WithCharset("UTF-8"), iris.WithoutVersionChecker, iris.WithoutServerError(iris.ErrServerClosed))
}

func info(ctx iris.Context) {
	method := ctx.Method()       // the http method requested a server's resource.
	subdomain := ctx.Subdomain() // the subdomain, if any.

	// the request path (without scheme and host).
	path := ctx.Path()
	// how to get all parameters, if we don't know
	// the names:
	paramsLen := ctx.Params().Len()

	ctx.Params().Visit(func(name string, value string) {
		ctx.Writef("%s = %s\n", name, value)
	})
	ctx.Writef("\nInfo\n\n")
	ctx.Writef("Method: %s\nSubdomain: %s\nPath: %s\nParameters length: %d", method, subdomain, path, paramsLen)
}
