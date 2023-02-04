package main

import (
	// "net/http"
	// "crypto/tls"
	// "golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
)

func main() {
	StartServerLogging()
	log.Println("Starting echo")
	e := echo.New()
	e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
	e.Use(middleware.Recover())
	
	e.Use(middleware.CORS())

	e.GET("/isloggedin", IsLoggedInHandler)

	e.GET("/new", CreateNewCommentHandler)

	e.GET("/all", GetAllCommentsHandler)

	e.GET("/accepted", GetAllAcceptedCommentsHandler)
	e.GET("/rejected", GetAllRejectedCommentsHandler)
	e.GET("/jailed", GetAllJailedCommentsHandler)

	e.GET("/accept", AcceptCommentHandler)
	// e.GET("/reject", RejectCommentHandler)

	e.GET("/newestreq", NewEstReqHandler)
	e.GET("/allest", GetAllEstimatesHandler)
	e.GET("/completeestreq", CompleteEstReqHandler)

	e.GET("/signup", SignUpHandler)
	e.GET("/signin", SignInHandler)
	e.GET("/signout", SignOutHandler)

	e.GET("/adminsignin", AdminSignInHandler)
	e.GET("/adminsignout", AdminSignOutHandler)

	e.Static("/static", "static") //for pics
	e.Logger.Fatal(e.StartAutoTLS(":443"))
	//e.Logger.Fatal(e.Start(":9090"))
}
