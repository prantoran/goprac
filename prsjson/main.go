package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type Req struct {
	RideID     string  `json:"ride_id"`
	FinalDist  float64 `json:"final_distance"`
	MethodUsed string  `json:"method_used"`
	CalcTime   float64 `json:"calculation_time"`
	Cohort     int     `json:"cohort"`
}

type SqlRet struct {
	ID          int     `json:"id"`
	HashedID    string  `json:"hashed_id"`
	RideDist    float64 `json:"ride_distance"`
	DistFromGPS float64 `json:"distance_from_gps"`
}

type Response struct {
	RideID      string  `json:"ride_id"`
	RideDist    float64 `json:"ride_distance"`
	DistFromGPS float64 `json:"distance_from_gps"`
	TomTomDist  float64 `json:"tomtom_distance"`
	DiffRDTom   float64 `json:"diff_rd_tom"`
	DiffGPSTom  float64 `json:"diff_gps_tom"`
}

func main() {
	// Open our jsonFile
	jsonFile, err := os.Open("rets.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened rets.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var reqs []Req

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &reqs)

	ids := ""

	reqCache := map[string]Req{}

	for i, r := range reqs {
		// fmt.Println(r)
		if i > 0 {
			ids += ","
		}

		ids += r.RideID

		reqCache[r.RideID] = r

	}

	fmt.Println("no. of ids:", len(reqs))

	sqlQuery := "select id, hashed_id, ride_distance, distance_from_gps from rides where id in (" + ids + ")"

	fmt.Println("sqlQuryL:", sqlQuery)

	// Open our jsonFile
	sqlResp, err := os.Open("sqlresp.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened rets.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer sqlResp.Close()

	// read our opened xmlFile as a byte array.
	byteValueSQL, _ := ioutil.ReadAll(sqlResp)
	fmt.Println("byteValueSQL:", string(byteValueSQL))
	// we initialize our Users array
	var sqlRet []SqlRet

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	err = json.Unmarshal(byteValueSQL, &sqlRet)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("no of sql resp:", len(sqlRet))

	resps := []Response{}

	for _, s := range sqlRet {
		fmt.Println(s)
		r := Response{}
		rideID := strconv.Itoa(s.ID)

		req := reqCache[rideID]

		r.RideID = rideID
		r.TomTomDist = req.FinalDist * 1000
		r.DistFromGPS = s.DistFromGPS
		r.RideDist = s.RideDist
		r.DiffRDTom = s.RideDist - req.FinalDist*1000
		r.DiffGPSTom = s.DistFromGPS - req.FinalDist*1000

		resps = append(resps, r)
	}

	respsJSON, _ := json.Marshal(resps)
	err = ioutil.WriteFile("difftom.json", respsJSON, 0644)
}
