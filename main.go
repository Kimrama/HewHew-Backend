package main

import (
	"encoding/json"
	"fmt"
	"hewhew-backend/config"
	"hewhew-backend/database"
	"hewhew-backend/server"
	//"hewhew-backend/database/migration"
)

func main() {
	conf := config.LoadConfig()
	confString, _ := json.MarshalIndent(conf, "", "  ")
	fmt.Println("config:\n", string(confString))

	database := database.NewPostgresDatabase(conf.Database)
	
	//migration.Migrate(database)
	server := server.NewFiberServer(conf, database)

	server.Start()

}
