package main

import (
	"encoding/json"
	"fmt"
	"hewhew-backend/config"
	"hewhew-backend/database"
)

func main() {
	conf := config.LoadConfig()
	confString, _ := json.MarshalIndent(conf, "", "  ")
	fmt.Println("config:\n", string(confString))

	database.NewPostgresDatabase(conf.Database)

}
