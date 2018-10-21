package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

type Refund struct {
	ID          string
	status      string
	Requester   string
	RequestDate string
	CityID      string
	RideID      string // hashID
	PhoneNumber string
}

func main() {
	filename := "refundclaim.csv"

	// Open CSV file
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		panic(err)
	}

	hashIDs := ""
	// Loop through lines & turn into object
	for i, line := range lines {
		if i > 40 {
			break
		}

		if i == 0 {
			continue
		}

		data := Refund{
			ID:     line[0],
			RideID: line[5],
		}

		if i > 1 {
			hashIDs += ","
		}

		fmt.Println(data.ID + " " + data.RideID)
		hashIDs += "\"" + data.RideID + "\""

	}

	// fmt.Println("hashIDs:", hashIDs)

	query := "select id, hashed_id, ride_distance, distance_from_gps from rides where hashed_id in (" + hashIDs + ")"

	fmt.Println("query:", query)
}
