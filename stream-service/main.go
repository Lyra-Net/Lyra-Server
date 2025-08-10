package main

import (
	"log"
	"net/http"
	"stream-service/config"
	"stream-service/routers"
)

func main() {
	r := routers.InitRouter()
	appEnv := config.InitConfig()
	log.Println("Stream service is running on: " + appEnv.PORT)
	if err := http.ListenAndServe(":"+appEnv.PORT, r); err != nil {
		log.Fatal(err)
	}
}
