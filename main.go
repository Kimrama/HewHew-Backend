package main

import (
	"encoding/json"
	"fmt"
	"hewhew-backend/config"
)

func main() {
	conf := config.LoadConfig()
	confString, _ := json.MarshalIndent(conf, "", "  ")
	fmt.Println("config:\n", string(confString))
}
