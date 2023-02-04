package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type EstimateS struct {
	UUID           string
	FileName       string
	DateTime       string
	Status         string
	Name           string
	Address        string
	City           string
	Telephone      string
	Email          string
	ReqServiceDate string
	Comment        string
	Photo          string
}

func NewEstReqHandler(c echo.Context) error {
	query := c.QueryParam("est")
	sp := strings.Split(query, "/")
	uuid, err := UUID()
	CheckError(err, "UUID creation failed")
	newdest := "./data/estimates/estimate_" + uuid + ".json"
	var est EstimateS
	est.UUID = uuid
	est.FileName = newdest
	est.DateTime = time.Now().String()
	est.Status = "active"
	est.Name = sp[0]
	est.Address = sp[1]
	est.City = sp[2]
	est.Telephone = sp[3]
	est.Email = sp[4]
	est.ReqServiceDate = sp[5]
	est.Comment = sp[6]
	est.Photo = sp[7]
	jsonEstimate, err := json.Marshal(est)
	CheckError(err, "JsonEstimate has failed")
	Writefile(jsonEstimate, newdest)
	return c.JSON(http.StatusOK, "Estimate Request Created")
}

func globEstimate() []EstimateS {
	accept, err := filepath.Glob("./data/estimates/*.json")
	CheckError(err, "Accept glob has failed")
	var allEstimate2 []EstimateS
	var allEstimate EstimateS
	for _, path := range accept {
		f, err := ioutil.ReadFile(path)
		CheckError(err, "Open has failed")
		json.Unmarshal(f, &allEstimate)
		allEstimate2 = append(allEstimate2, allEstimate)
	}
	return allEstimate2
}

func readEstFile(afile string) EstimateS {
	estimate, err := os.ReadFile(afile)
	CheckError(err, "Read file has failed")
	var est EstimateS
	json.Unmarshal(estimate, &est)
	return est
}

func CompleteEstReqHandler(c echo.Context) error {
	query := c.QueryParam("uuid")
	globestimates, err := filepath.Glob("./data/estimates/*.json")
	CheckError(err, "completed glob has failed")
	var AEst EstimateS
	for _, f := range globestimates {
		if strings.Contains(query, f) {
			esti := readEstFile(f)
			newdest := "./data/estcompleted/estcompleted_" + esti.UUID + ".json"
			oldFile := esti.FileName
			AEst.UUID = esti.UUID
			AEst.FileName = newdest
			AEst.DateTime = esti.DateTime
			AEst.Status = "completed"
			AEst.Name = esti.Name
			AEst.Address = esti.Email
			AEst.City = esti.City
			AEst.Telephone = esti.Telephone
			AEst.Email = esti.Email
			AEst.ReqServiceDate = esti.ReqServiceDate
			AEst.Comment = esti.Comment
			AEst.Photo = esti.Photo
			jsonCompleted, err := json.Marshal(AEst)
			CheckError(err, "json marshal has failed")
			Writefile(jsonCompleted, newdest)
			os.Remove(oldFile)
		}
	}
	return c.JSON(http.StatusOK, "Estimate Request Created")
}

func GetAllEstimatesHandler(c echo.Context) error {
	esti := globEstimate()
	if len(esti) != 0 {
		return c.JSON(http.StatusOK, esti)
	} else {
		return c.JSON(http.StatusOK, "None")
	}
}
