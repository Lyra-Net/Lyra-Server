package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type APPENV struct {
	PORT string
}

func InitConfig() APPENV {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No Env found")
	}
	Port := os.Getenv("PORT")
	return APPENV{PORT: Port}
}
