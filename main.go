package main

import (
	"database/sql"
	"log"

	"github.com/Ecc-asplay/backend/util"
)

func main() {
	config, err := util.LoadConfig("./")
	if err != nil {
		log.Println("app.env 見つけてない")
	}

	psql, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Println("DB 接続できない")
	}

	if err := psql.Ping(); err != nil {
		log.Println("Error ：", err.Error())
		return
	}

}
