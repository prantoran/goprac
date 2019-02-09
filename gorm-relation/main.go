package main

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Person struct {
	PersonID  int32  `gorm:"column:id;primary_key"`
	LastName  string `gorm:"column:last_name"`
	FirstName string `gorm:"column:first_name"`
	Adress    string `gorm:"column:address"`
	City      string `gorm:"column:city"`
}

func (p Person) TableName() string {
	return "Persons"
}

type Order struct {
	OrderID     int32  `gorm:"column:id;primary_key"`
	OrderNumber int32  `gorm:"column:order_number"`
	PersonIDs   int32  `gorm:"column:person_id"`
	Persons     Person `gorm:"foreignkey:PersonIDs"`
}

func (o Order) TableName() string {
	return "Orders"
}

func chkPersonsTable(db *gorm.DB) {
	if db.HasTable("Persons") {
		log.Println("Persons table exists")
	}

	if db.HasTable(&Person{}) {
		log.Println("Person{} table exists")
	}

	log.Println(db.CreateTable(&Person{}).Error)

	persons := []Person{}

	err := db.Find(&persons).Error
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	for _, u := range persons {
		log.Println("person:", u)
	}
}

func chkOrdersTable(db *gorm.DB) {
	if db.HasTable("Orders") {
		log.Println("Orders table exists")
	}

	if db.HasTable(&Order{}) {
		log.Println("Order{} table exists")
	}

	log.Println(db.CreateTable(&Order{}).Error)

	order := Order{}

	err := db.Preload("Persons").Find(&order).Error

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Println("order: personIDs", order.PersonIDs, " orderID:", order.OrderID)
	log.Println("order:", order)

}

func main() {
	db, err := gorm.Open("mysql", "root:example@tcp(127.0.0.1:3306)/location?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()

	db = db.Debug()

	if err != nil {
		log.Fatal(err)
	}

	chkPersonsTable(db)

	chkOrdersTable(db)

}
