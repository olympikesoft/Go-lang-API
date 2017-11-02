package database

import (
	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"config"
)

var Mysql *sqlx.DB

func Connect() {
	var err error

	if Mysql, err = sqlx.Connect("mysql",  config.GetConfig().DbDsn); err != nil {
		logrus.Fatalf("Connection Error", err)
	}

	// Check if is alive
	if err = Mysql.Ping(); err != nil {
		logrus.Fatalf("Ping Error", err)
	}
}
