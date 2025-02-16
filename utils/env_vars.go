package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvVariables struct {
	EmailAppPassword string
	AppEmailID string
	DEVDATABASEURL string
}

func GetEnvVars() (EnvVariables, error) {
	var e EnvVariables

	err := godotenv.Load(".env")
	if err !=  nil{
		log.Fatal(err)
		return e,err
	}

	e.EmailAppPassword = os.Getenv("EMAILAPPPASSWORD")
	e.AppEmailID = os.Getenv("APPEMAILID")
	e.DEVDATABASEURL = os.Getenv("DEVDBURL")

	return e,nil
}