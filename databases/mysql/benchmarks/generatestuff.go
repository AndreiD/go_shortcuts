package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"squsers/database"
	"squsers/models"
)

func main() {
	fmt.Println("Starting filling up the database....")

	database.InitMySQLDatabase()

	pwd, _ := os.Getwd()
	jsonFile, err := os.Open(pwd + "/mock_users.json")
	if err != nil {
		log.Fatalf("cant find the file %s", err.Error())
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var users []models.User
	err = json.Unmarshal(byteValue, &users)
	if err != nil {
		log.Fatalf("error unmarshaling %s", err.Error())
	}

	// for testing
	var aUser models.User
	aUser.ID = "1a2b4c"
	aUser.Role = "user"
	aUser.FirstName = "XXXX"
	aUser.LastName = "YYYYYY"
	_ = database.CreateUser(aUser)

	for _, user := range users {
		_ = database.CreateUser(user)
	}

	fmt.Println("------------ done.")

}
