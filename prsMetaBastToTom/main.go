package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type LatLons struct {
	Lats []float64 `json:"lats"`
	Lons []float64 `json:"longs"`
}

type Req struct {
	RideID   string    `json:"ride_id"`
	CityID   int32     `json:"city_id"`
	RideType int32     `json:"ride_type"`
	Mode     string    `json:"mode"`
	Data     []LatLons `json:"data"`
	EstDis   float64   `json:"estimated_distance"`
	EstDur   float64   `json:"estimated_duration"`
}

type Resp struct {
	RideID    string  `json:"ride_id"`
	FinalDist float64 `json:"final_distance"`
	Method    string  `json:"method_used"`
	CalcTm    float64 `json:"calculation_time"`
	Cohort    int     `json:"cohort"`
}

func main() {
	csvFile, err := os.Open("./d.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1

	csvData, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// var emp Employee
	// var employees []Employee
	//

	cache := map[string]Req{}

	for _, each := range csvData {
		// fmt.Printf("id: %v\trideID: %v type: %T\ttype: %v\tdata: %v type: %T\tcreatedAt: %v type: %T\n", each[0], each[1], each[1], each[2], each[3], each[3], each[4], each[4])

		rideID := each[1]

		var r Req
		var ok bool

		if r, ok = cache[rideID]; !ok {
			r = Req{}
			latlon := LatLons{}
			r.Data = append(r.Data, latlon)
			r.CityID = 2
			r.RideType = 1
			r.RideID = rideID
		}

		rawlatlon := each[3]

		lltoks := strings.Split(rawlatlon, ",")

		if len(lltoks) < 2 {
			continue
		}
		// fmt.Println("lltoks:", lltoks)
		lat, err := strconv.ParseFloat(lltoks[0], 64)
		// fmt.Println("lat:", lat)

		if err != nil {
			log.Println(err)
			continue
		}

		// fmt.Println("lat:", lat)

		lon, err := strconv.ParseFloat(lltoks[1], 64)
		if err != nil {
			log.Println(err)
			continue
		}
		// fmt.Println("lon:", lon)

		r.Data[0].Lats = append(r.Data[0].Lats, lat)
		r.Data[0].Lons = append(r.Data[0].Lons, lon)

		cache[rideID] = r

		// fmt.Println("req:", r)
	}

	rets := []Resp{}
	url := "http://127.0.0.1:8000/api/tomtom/calculate_distance"
	for k, v := range cache {

		pings := len(v.Data[0].Lats)

		v.EstDis = (float64(pings) * 5 * 16) / (60 * 60)
		v.EstDur = float64(pings) * 5

		cache[k] = v

		// fmt.Println("key:", k, " pings:", pings, " estdur:", cache[k].EstDur, " estdis:", cache[k].EstDis)
		b, err := json.Marshal(v)
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}
		// fmt.Println("reqbody:", string(b))
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))

		if err != nil {
			log.Println("new req:", err)
		}
		// req.Header.Set("X-Custom-Header", "myvalue")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		// reqs = append(reqs, v)
		ret := Resp{}

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&ret)
		if err != nil {
			panic(err)
		}
		log.Println(ret)

		time.Sleep(500 * time.Microsecond)

		ret.RideID = v.RideID

		rets = append(rets, ret)

		resp.Body.Close()

	}

	// Convert to JSON
	// jsonData, err := json.Marshal(reqs)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// fmt.Println(string(jsonData))

	// jsonFile, err := os.Create("./data.json")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer jsonFile.Close()

	// jsonFile.Write(jsonData)
	// // jsonFile.Close()

	// Convert to JSON
	jsonData, err := json.Marshal(rets)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))

	jsonFile, err := os.Create("./rets.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	jsonFile.Write(jsonData)
	// jsonFile.Close()
}
