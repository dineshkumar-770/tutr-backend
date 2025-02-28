package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvVariables struct {
	EmailAppPassword  string
	AppEmailID        string
	DatabaseRole      string
	DatabasePassword  string
	DatabaseName      string
	DatabasePORT      string
	DatabaseIPAddress string
}

func GetEnvVars() (EnvVariables, error) {
	var e EnvVariables

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
		return e, err
	}

	e.EmailAppPassword = os.Getenv("EMAILAPPPASSWORD")
	e.AppEmailID = os.Getenv("APPEMAILID")
	e.DatabaseRole = os.Getenv("DEVDBROLE")
	e.DatabasePassword = os.Getenv("DEVDBPASSWORD")
	e.DatabaseIPAddress = os.Getenv("DEVDBIPADDRESS")
	e.DatabasePORT = os.Getenv("DATABASEPORT")
	e.DatabaseName = os.Getenv("DEVDBNAME")

	return e, nil
}
