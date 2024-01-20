package main

import (
	"fmt"
	"os"

	"github.com/NATCHAYATP/E-Commerce/config"
	"github.com/NATCHAYATP/E-Commerce/modules/servers"
	"github.com/Rayato159/kawaii-shop-tutorial/pkg/databases"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func main() {
	cfg := config.LoadConfig(envPath())

	db := databases.DbConnect(cfg.Db())
	defer db.Close()
	fmt.Println(db)

	servers.NewServer(cfg, db).Start()
}
