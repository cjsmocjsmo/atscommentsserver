package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"github.com/labstack/echo/v4"
)

type ProfileS struct {
	AccountID string
	UserID    string
	Email     string
	Password  string
	OldFile   string
}

type LoggedIn struct {
	User string
}

func SignUpHandler(c echo.Context) error {
	query := c.QueryParam("creds")
	sp := strings.Split(query, "/")
	acid, err := UUID()
	CheckError(err, "acid #1 has failed")
	userid := acid + "_" + sp[0]
	newdest := "./data/admin/profiles/" + acid + "_" + sp[0] + ".json"
	var pfile ProfileS
	pfile.Email = sp[0]
	pfile.Password = sp[1]
	pfile.AccountID = acid
	pfile.UserID = userid
	pfile.OldFile = newdest
	jsonProfile, err := json.Marshal(pfile)
	CheckError(err, "Profile creation has failed")
	Writefile(jsonProfile, newdest)
	return c.JSON(http.StatusOK, userid)
}

func addToLoggedInList(cstring string) string {
	var newuser LoggedIn
	newuser.User = cstring
	loc := "./data/admin/loggedInList.json"
	uil, err := os.ReadFile(loc)
	CheckError(err, "addToLoggedInList read file has failed")
	var loggedinlist []LoggedIn
	json.Unmarshal(uil, &loggedinlist)
	loggedinlist = append(loggedinlist, newuser)
	jsonloggedinlist, errr := json.Marshal(loggedinlist)
	CheckError(errr, "jsonloggedinlist has failed")
	Writefile(jsonloggedinlist, loc)
	return "User is logged in"
}

func checkForAccount(astring string) bool {
	profiles, err := filepath.Glob("./data/admin/profiles/*.json")
	CheckError(err, "Glob profiles has failsed")
	result := false
	for _, pf := range profiles {
		log.Println(pf)
		sp1 := filepath.Base(pf)
		idx := len(sp1) - 5
		sp := sp1[:idx]
		log.Println(sp)
		log.Println(astring)
		if sp != astring {
			result = false
		} else {
			result = true
		}
	}
	return result
}

func checkForAlreadyLoggedIn(creds string) bool {
	// for internal use
	loc := "./data/admin/loggedInList.json"
	lil, err := os.ReadFile(loc)
	CheckError(err, "checkForAlreadyLoggedIn readfile has failed")
	var logedinlist []LoggedIn
	json.Unmarshal(lil, &logedinlist)
	result := false
	if len(logedinlist) != 0 {
		for _, l := range logedinlist {
			if l.User == creds {
				result = true
			}
		}
	}
	return result
}

func IsLoggedInHandler(c echo.Context) error {
	// for external use
	query := c.QueryParam("creds")
	loc := "./data/admin/loggedInList.json"
	lil, err := os.ReadFile(loc)
	CheckError(err, "checkForAlreadyLoggedIn readfile has failed")
	var logedinlist []LoggedIn
	json.Unmarshal(lil, &logedinlist)
	result := false
	if len(logedinlist) != 0 {
		for _, l := range logedinlist {
			if l.User == query {
				result = true
			}
		}
	}
	return c.JSON(http.StatusOK, result)
}

func SignInHandler(c echo.Context) error {
	log.Println("Starting SignInHandler")
	query := c.QueryParam("creds")
	// sp := strings.Split(query, "/")
	sp := strings.Replace(query, "/", "_", 1)
	hasAccount := checkForAccount(sp)
	log.Println(hasAccount)

	if hasAccount {
		// creds := strings.Replace(query, "/", "_", 1)
		// check to see of already logged in
		if checkForAlreadyLoggedIn(sp) {
			return c.JSON(http.StatusOK, "User is already logged in")
		} else {
			results := addToLoggedInList(sp)
			log.Println("SignInHandler complete")
			return c.JSON(http.StatusOK, results)
		}
	} else {
		log.Println("SignInHandler complete")
		return c.JSON(http.StatusOK, "Please Create An Account")
	}

}

func getIndex(query string, loggedinlist []LoggedIn) int {
	idxx := 0
	for idx, UID := range loggedinlist {
		if UID.User == query {
			idxx = idx
		}
	}
	return idxx
}

func SignOutHandler(c echo.Context) error {
	query := c.QueryParam("creds")
	loc := "./data/admin/loggedInList.json"
	lil, err := os.ReadFile(loc)
	CheckError(err, "Signout has failed")
	var loggedinlist []LoggedIn
	json.Unmarshal(lil, &loggedinlist)
	idx := getIndex(query, loggedinlist)
	loggedinlist = append(loggedinlist[:idx], loggedinlist[idx+1:]...)
	jsonloggedinlist, err := json.Marshal(loggedinlist)
	CheckError(err, "jsonloggedinlist marshal has failed")
	Writefile(jsonloggedinlist, loc)
	return c.JSON(http.StatusOK, "User is logged out")
}

func AdminSignInHandler(c echo.Context) error {
	query := c.QueryParam("comment")
	// sp := strings.Split(query, "/")

	return c.JSON(http.StatusOK, query)
}

func AdminSignOutHandler(c echo.Context) error {
	query := c.QueryParam("comment")
	// sp := strings.Split(query, "/")

	return c.JSON(http.StatusOK, query)
}
