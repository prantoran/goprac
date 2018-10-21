package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/tealeg/xlsx"
)

const (
	MedianRolling string = "MEDIAN ROLLING"
)

type Resp struct {
	RideID    string  `json:"ride_id"`
	FinalDist float64 `json:"final_distance,omitempty"`
	Method    string  `json:"method_used,omitempty"`
	CalcTm    float64 `json:"calculation_time,omitempty"`
	Cohort    int     `json:"cohort,omitempty"`
	Err       string  `json:"error,omitempty"`
}

type SQLResp struct {
	ID          int     `json:"id"`
	HashedID    string  `json:"hashed_id"`
	RideDist    float64 `json:"ride_distance"`
	EstRideDist float64 `json:"estimated_ride_distance"`
	DistFromGPS float64 `json:"distance_from_gps"`
}

type CmpResp struct {
	HashedID    string  `json:"hashed_id"`
	RideID      int     `json:"ride_id"`
	TomPyDist   float64 `json:"tompy_dist"`
	TomGoDist   float64 `json:"tomgo_dist"`
	MethodGo    string  `json:"method_used_go,omitempty"`
	MethodPy    string  `json:"method_used_py,omitempty"`
	RideDist    float64 `json:"ride_distance"`
	EstRideDist float64 `json:"estimated_ride_distance,omitempty"`
	DistFromGPS float64 `json:"dist_from_gps"`
	DifSQLGo    float64 `json:"dif_sql_go"`
	DifSQLPy    float64 `json:"dif_sql_py"`
	DifEstSQLGo float64 `json:"dif_estsql_go"`
	DifEstSQLPy float64 `json:"dif_estsql_py"`
	DifGoPy     float64 `json:"dif_go_py"`
	DifGPSGo    float64 `json:"dif_gps_go"`
	DifGPSPy    float64 `json:"dif_gps_py"`
}

func main() {
	gotom := respFrmFile("./retsgo.json")
	pytom := respFrmFile("./retspy.json")

	rideq := respSQLFrmFile("./ridesq.json")
	// fmt.Println(rideq)

	cacheSQL := map[int]SQLResp{}
	for _, s := range rideq {
		cacheSQL[s.ID] = s
	}

	cacheGo := map[int]Resp{}
	for _, r := range gotom {
		id, err := strconv.Atoi(r.RideID)
		if err != nil {
			logrus.Warn("gotom RideID:", r.RideID, " err:", err)
			continue
		}

		cacheGo[id] = r
	}

	cachePy := map[int]Resp{}
	for _, r := range pytom {
		id, err := strconv.Atoi(r.RideID)
		if err != nil {
			logrus.Warn("pytom RideID:", r.RideID, " err:", err)
			continue
		}
		r.FinalDist *= 1000

		cachePy[id] = r
	}

	cmpResps := []CmpResp{}

	for _, s := range rideq {
		id := s.ID
		if cacheGo[id].Method != MedianRolling || cachePy[id].Method != MedianRolling {
			continue
		}
		u := CmpResp{}
		u.RideID = id
		u.HashedID = s.HashedID
		u.RideDist = s.RideDist
		u.EstRideDist = s.EstRideDist
		u.TomGoDist = cacheGo[id].FinalDist
		u.TomPyDist = cachePy[id].FinalDist
		u.MethodGo = cacheGo[id].Method
		u.MethodPy = cachePy[id].Method
		u.DifSQLGo = u.RideDist - u.TomGoDist
		u.DifSQLPy = u.RideDist - u.TomPyDist
		u.DifGoPy = u.TomGoDist - u.TomPyDist
		u.DifEstSQLGo = u.EstRideDist - u.TomGoDist
		u.DifEstSQLPy = u.EstRideDist - u.TomPyDist
		u.DistFromGPS = s.DistFromGPS
		u.DifGPSPy = u.DistFromGPS - u.TomPyDist
		u.DifGPSGo = u.DistFromGPS - u.TomGoDist
		cmpResps = append(cmpResps, u)
	}

	writeFile("./sakura.json", cmpResps)

	writeExcel("./tomtomcmp.xlsx", cmpResps)
}

func addRowExcel(sheet *xlsx.Sheet, fields ...string) {

	var row *xlsx.Row
	var cell *xlsx.Cell

	row = sheet.AddRow()

	for _, s := range fields {
		cell = row.AddCell()
		cell.Value = s
	}
}

func writeExcel(fname string, rets []CmpResp) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("stats")
	if err != nil {
		fmt.Printf(err.Error())
	}

	addRowExcel(sheet, "HashedID", "TomPy", "TomGo", "RideDist", "EstRideDist", "DistFromGPS")

	for i, v := range rets {

		fmt.Println("i:", i, "\tv:", v, "\tgps:", v.DistFromGPS)

		if v.MethodGo != MedianRolling || v.MethodPy != MedianRolling {
			continue
		}

		addRowExcel(
			sheet,
			v.HashedID,
			strconv.Itoa(int(v.TomPyDist)),
			strconv.Itoa(int(v.TomGoDist)),
			fmt.Sprintf("%v", v.RideDist),
			fmt.Sprintf("%v", v.EstRideDist),
			fmt.Sprintf("%v", v.DistFromGPS))
	}

	err = file.Save(fname)
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func writeFile(fname string, rets []CmpResp) {
	// Convert to JSON
	jsonData, err := json.Marshal(rets)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// fmt.Println(string(jsonData))

	jsonFile, err := os.Create(fname)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	jsonFile.Write(jsonData)
	// jsonFile.Close()
}

func respFrmFile(fname string) []Resp {
	jsonFile, err := os.Open(fname)
	if err != nil {
		logrus.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var resps []Resp
	json.Unmarshal(byteValue, &resps)

	return resps
}

func respSQLFrmFile(fname string) []SQLResp {
	jsonFile, err := os.Open(fname)
	if err != nil {
		logrus.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var resps []SQLResp
	json.Unmarshal(byteValue, &resps)
	return resps
}
