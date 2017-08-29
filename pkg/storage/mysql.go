package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	// driver to specifically connect to mysql
	_ "github.com/go-sql-driver/mysql"
)

// MySQLDBConf is a conf for the mysql database
type MySQLDBConf struct {
	Protocol string
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
}

// NewMySQLDBConnPool connects to db and return a connection pool
func NewMySQLDBConnPool(mysqlDBConf *MySQLDBConf) (*sqlx.DB, error) {
	dsn := mysqlDBConf.User + ":" +
		mysqlDBConf.Password + "@" +
		mysqlDBConf.Protocol + "(" +
		mysqlDBConf.Host + ":" +
		mysqlDBConf.Port + ")/" +
		mysqlDBConf.DbName + "?parseTime=true&multiStatements=true"

	pool, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("NewMySQLDBConnPool: sqlx.Open %v", err)
	}

	errP := pool.Ping()
	if errP != nil {
		return nil, fmt.Errorf("NewMySQLDBConnPool: pool.Ping %v", errP)
	}

	return pool, nil
}
