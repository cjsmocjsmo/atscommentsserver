package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type CommentS struct {
	Name       string
	Email      string
	Comment    string
	UUID       string
	DateTime   string
	Accepted   string
	Rejected   string
	Jailed     string
	StarRating string
	FileName   string
}

func StartServerLogging() string {
	logtxtfile := "./log/logfile.txt"
	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile(logtxtfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println(err)
	}
	log.SetOutput(file)
	fmt.Println("Logging started")
	return "Server logging started"
}

func UUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := rand.Read(uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	uuid[8] = 0x80
	uuid[4] = 0x40
	boo := hex.EncodeToString(uuid)
	return boo, nil
}

func CheckError(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		fmt.Println(err)
		log.Println(msg)
		log.Println(err)
		panic(err)
	}
}

func Readfile(afile string) CommentS {
	comment, err := os.ReadFile(afile)
	CheckError(err, "Read file has failed")
	var coms CommentS
	json.Unmarshal(comment, &coms)
	return coms
}

func Writefile(com []byte, dest string) {
	f, err := os.Create(dest)
	CheckError(err, "Creating file has failed")
	defer f.Close()
	_, err = f.Write(com)
	CheckError(err, "Write has failed")
}

func CreateNewCommentHandler(c echo.Context) error {
	uuid, err := UUID()
	newdest := "./data/jailed/jailed_" + uuid + ".json"
	CheckError(err, "uuid creation has failed")
	// query string needs to be in the format
	// ?comment=John Doe/john@gmail/Job well done
	query := c.QueryParam("comment")
	sp := strings.Split(query, "/")
	var NewCom CommentS
	NewCom.Name = sp[0]
	NewCom.Email = sp[1]
	NewCom.Comment = sp[2]
	NewCom.UUID = uuid
	NewCom.DateTime = time.Now().String()
	NewCom.Accepted = "No"
	NewCom.Rejected = "No"
	NewCom.Jailed = "yes"
	NewCom.StarRating = "5"
	NewCom.FileName = newdest
	jsonComment, err := json.Marshal(NewCom)
	CheckError(err, "Json marshal has failed")
	Writefile(jsonComment, newdest)
	return c.JSON(http.StatusOK, "Comment Created")
}

func globAccepted() []CommentS {
	accept, err := filepath.Glob("./data/accepted/*.json")
	CheckError(err, "Accept glob has failed")
	var allAccepted2 []CommentS
	var allAccepted CommentS
	for _, path := range accept {
		f, err := ioutil.ReadFile(path)
		CheckError(err, "Open has failed")
		json.Unmarshal(f, &allAccepted)
		allAccepted2 = append(allAccepted2, allAccepted)
	}
	return allAccepted2
}

func globRejected() []CommentS {
	rejected, err := filepath.Glob("./data/rejected/*.json")
	CheckError(err, "Rejected glob has failed")
	var allRejected2 []CommentS
	var allRejected CommentS
	for _, path := range rejected {
		f, err := ioutil.ReadFile(path)
		CheckError(err, "Open has failed")
		json.Unmarshal(f, &allRejected)
		allRejected2 = append(allRejected2, allRejected)
	}
	return allRejected2
}

func globJailed() []CommentS {
	jailed, err := filepath.Glob("./data/jailed/*.json")
	CheckError(err, "Rejected glob has failed")
	var allJailed CommentS
	var allJailed2 []CommentS
	for _, path := range jailed {
		f, err := ioutil.ReadFile(path)
		CheckError(err, "Open has failed")
		json.Unmarshal(f, &allJailed)
		allJailed2 = append(allJailed2, allJailed)
	}
	return allJailed2
}

func GetAllCommentsHandler(c echo.Context) error {
	allAccepted := globAccepted()
	allRejected := globRejected()
	allJailed := globJailed()
	var allList []CommentS
	if len(allAccepted) != 0 {
		allList = append(allList, allAccepted...)
	}
	if len(allRejected) != 0 {
		allList = append(allList, allRejected...)
	}
	if len(allJailed) != 0 {
		allList = append(allList, allJailed...)
	}
	// jsonAllList, err := json.Marshal(allList)
	// CheckError(err, "JsonallList has failed")

	return c.JSON(http.StatusOK, allList)
}

func GetAllAcceptedCommentsHandler(c echo.Context) error {
	globaccepted := globAccepted()
	if len(globaccepted) != 0 {
		var all []CommentS
		all = append(all, globaccepted...)
		return c.JSON(http.StatusOK, all)
	} else {
		return c.JSON(http.StatusOK, "None")
	}
}

func GetAllRejectedCommentsHandler(c echo.Context) error {
	globrejected := globRejected()
	if len(globrejected) != 0 {
		var all []CommentS
		all = append(all, globrejected...)
		return c.JSON(http.StatusOK, all)
	} else {
		return c.JSON(http.StatusOK, "None")
	}
}

func GetAllJailedCommentsHandler(c echo.Context) error {
	globjailed := globJailed()
	if len(globjailed) != 0 {
		var all []CommentS
		all = append(all, globjailed...)
		return c.JSON(http.StatusOK, all)
	} else {
		return c.JSON(http.StatusOK, "None")
	}
}

func AcceptCommentHandler(c echo.Context) error {
	query := c.QueryParam("uuid")
	globjailed, err := filepath.Glob("./data/jailed/*.json")
	CheckError(err, "globjail has failed")
	var ACom CommentS
	for _, f := range globjailed {
		if strings.Contains(query, f) {
			comment := Readfile(f)
			oldFile := comment.FileName
			dest := "/accepted/accepted_" + comment.UUID + ".json"
			ACom.Name = comment.Name
			ACom.Email = comment.Email
			ACom.Comment = comment.Comment
			ACom.UUID = comment.UUID
			ACom.DateTime = comment.DateTime
			ACom.Accepted = "yes"
			ACom.Rejected = "no"
			ACom.Jailed = "no"
			ACom.StarRating = comment.StarRating
			ACom.FileName = dest
			acom, err := json.Marshal(ACom)
			CheckError(err, "acom has failed")
			Writefile(acom, dest)
			os.Remove(oldFile)
		}
	}
	return c.JSON(http.StatusOK, "Comment Accemped")
}

// func RejectCommentHandler(c echo.Context) error {
// 	query := c.QueryParam("uuid")
// 	globrejected, err := filepath.Glob("./data/rejected/*.json")
// 	CheckError(err, "globrejected has failed")

// 	return c.JSON(http.StatusOK, "fuck")
// }
