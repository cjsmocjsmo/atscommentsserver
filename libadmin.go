package main

import (
	"archive/zip"
	"encoding/json"
	"io"

	// "io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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

type CredS struct {
	Token    string
	Email    string
	Password string
}

func TestHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "Server Is Up and Running")
}

func SignUpHandler(c echo.Context) error {
	query := c.QueryParam("creds")
	sp := strings.Split(query, "/")
	acid, err := UUID()
	CheckError(err, "acid #1 has failed")
	userid := acid + "_" + sp[0]
	newdest := "/home/porthose_cjsmo_cjsmo/atscommentsserver/data/admin/profiles/" + acid + "_" + sp[0] + ".json"
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
	loc := "/home/porthose_cjsmo_cjsmo/atscommentsserver/data/admin/loggedInList.json"
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
	profiles, err := filepath.Glob("/home/porthose_cjsmo_cjsmo/atscommentsserver/data/admin/profiles/*.json")
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
	loc := "/home/porthose_cjsmo_cjsmo/atscommentsserver/data/admin/loggedInList.json"
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
	loc := "/home/porthose_cjsmo_cjsmo/atscommentsserver/data/admin/loggedInList.json"
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
	loc := "/home/porthose_cjsmo_cjsmo/atscommentsserver/data/admin/loggedInList.json"
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

func ReadCredsFile(afile string) CredS {
	creditials, err := os.ReadFile(afile)
	CheckError(err, "Read file has failed")
	var creds CredS
	json.Unmarshal(creditials, &creds)
	return creds
}

func checkCreds(e1 string, e2 string) bool {
	if e1 != e2 {
		return false
	} else {
		return true
	}
}

func checkResults(e1 bool, e2 bool, e3 bool) bool {
	if !e1 || !e2 || !e3 {
		return false
	} else {
		return true
	}
}

func AdminSignInHandler(c echo.Context) error {
	query := c.QueryParam("creds")
	sp := strings.Split(query, "/")
	token := sp[0]
	email := sp[1]
	pword := sp[2]
	acreds := ReadCredsFile("/home/porthose_cjsmo_cjsmo/atscommentsserver/creds/admin_creds.json")
	result1 := checkCreds(token, acreds.Token)
	result2 := checkCreds(email, acreds.Email)
	result3 := checkCreds(pword, acreds.Password)
	isLoggedIn := checkResults(result1, result2, result3)
	return c.JSON(http.StatusOK, isLoggedIn)
}

func AdminSignOutHandler(c echo.Context) error {
	query := c.QueryParam("creds")
	// sp := strings.Split(query, "/")

	return c.JSON(http.StatusOK, query)
}

// func GlobAdmin() []ProfileS {
// 	profiles, err := filepath.Glob("./data/admin/profiles/*.json")
// 	CheckError(err, "Admin glob has failed")
// 	var allProfiles2 []ProfileS
// 	var allProfiles ProfileS
// 	for _, path := range profiles {
// 		f, err := ioutil.ReadFile(path)
// 		CheckError(err, "Open has failed")
// 		json.Unmarshal(f, &allProfiles)
// 		allProfiles2 = append(allProfiles2, allProfiles)
// 	}
// 	return allProfiles2
// }

func Backup(c echo.Context) error {
	log.Println("creating zip archive...")

	dt := time.Now()
	date := dt.Format("13-01-2022")

	addr := "/home/porthose_cjsmo_cjsmo/atscommentsserver/static/" + date + "_backup.zip"
	archive, err := os.Create(addr)
	if err != nil {
		panic(err)
	}
	defer archive.Close()
	zipWriter := zip.NewWriter(archive)

	g1, err := filepath.Glob("/home/porthose_cjsmo_cjsmo/atscommentsserver/data/accepted/*.json")
	CheckError(err, "accept glob has failed")
	if len(g1) != 0 {
		for _, g := range g1 {
			log.Println("opening first file...")
			f1, err := os.Open(g)
			if err != nil {
				panic(err)
			}
			defer f1.Close()

			base := filepath.Base(g)

			log.Println("writing first file to archive...")
			foo := "/home/porthose_cjsmo_cjsmo/atscommentsserver/accepted/" + base
			w1, err := zipWriter.Create(foo)
			if err != nil {
				panic(err)
			}
			if _, err := io.Copy(w1, f1); err != nil {
				panic(err)
			}
		}
	}

	g2, err := filepath.Glob("/home/porthose_cjsmo_cjsmo/atscommentsserver/data/estimates/*.json")
	CheckError(err, "estimates glob has failed")
	if len(g2) != 0 {
		for _, g := range g1 {
			log.Println("opening first file...")
			f1, err := os.Open(g)
			if err != nil {
				panic(err)
			}
			defer f1.Close()

			base := filepath.Base(g)

			log.Println("writing first file to archive...")
			foo := "/home/porthose_cjsmo_cjsmo/atscommentsserver/estimates/" + base
			w1, err := zipWriter.Create(foo)
			if err != nil {
				panic(err)
			}
			if _, err := io.Copy(w1, f1); err != nil {
				panic(err)
			}
		}
	}

	g3, err := filepath.Glob("/home/porthose_cjsmo_cjsmo/atscommentsserver/data/jailed/*.json")
	CheckError(err, "jailed glob has failed")
	if len(g3) != 0 {
		for _, g := range g1 {
			log.Println("opening first file...")
			f1, err := os.Open(g)
			if err != nil {
				panic(err)
			}
			defer f1.Close()

			base := filepath.Base(g)

			log.Println("writing first file to archive...")
			foo := "/home/porthose_cjsmo_cjsmo/atscommentsserver/jailed/" + base
			w1, err := zipWriter.Create(foo)
			if err != nil {
				panic(err)
			}
			if _, err := io.Copy(w1, f1); err != nil {
				panic(err)
			}
		}
	}

	g4, err := filepath.Glob("/home/porthose_cjsmo_cjsmo/atscommentsserver/data/rejected/*.json")
	CheckError(err, "rejected glob failed")
	if len(g4) != 0 {
		for _, g := range g1 {
			log.Println("opening first file...")
			f1, err := os.Open(g)
			if err != nil {
				panic(err)
			}
			defer f1.Close()

			base := filepath.Base(g)

			log.Println("writing first file to archive...")
			foo := "/home/porthose_cjsmo_cjsmo/atscommentsserver/rejected/" + base
			w1, err := zipWriter.Create(foo)
			if err != nil {
				panic(err)
			}
			if _, err := io.Copy(w1, f1); err != nil {
				panic(err)
			}
		}
	}

	g6, err := filepath.Glob("/home/porthose_cjsmo_cjsmo/atscommentsserver/data/admin/profiles/*.json")
	CheckError(err, "admin glob has failed")
	if len(g6) != 0 {
		for _, g := range g1 {
			log.Println("opening first file...")
			f1, err := os.Open(g)
			if err != nil {
				panic(err)
			}
			defer f1.Close()

			base := filepath.Base(g)

			log.Println("writing first file to archive...")
			foo := "/home/porthose_cjsmo_cjsmo/atscommentsserver/admin/profiles/" + base
			w1, err := zipWriter.Create(foo)
			if err != nil {
				panic(err)
			}
			if _, err := io.Copy(w1, f1); err != nil {
				panic(err)
			}
		}
	}

	f1, err := os.Open("/home/porthose_cjsmo_cjsmo/atscommentsserver/data/admin/loggedInList.json")
	if err != nil {
		panic(err)
	}
	defer f1.Close()

	log.Println("writing first file to archive...")
	foo := "/home/porthose_cjsmo_cjsmo/atscommentsserver/admin/loggedInList.json"
	w1, err := zipWriter.Create(foo)
	if err != nil {
		panic(err)
	}
	if _, err := io.Copy(w1, f1); err != nil {
		panic(err)
	}

	log.Println("closing zip archive...")
	zipWriter.Close()

	return c.JSON(http.StatusOK, "backup complete")
}
