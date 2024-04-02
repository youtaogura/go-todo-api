package database

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

var conn *gorm.DB

func Setup() *gorm.DB {
	connection, err := newConnection(
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)
	if err != nil {
		fmt.Println("Failed to connect to database")
		fmt.Println(err)
		panic(err)
	}
	conn = connection
	return conn
}

func GetConnection() *gorm.DB {
	if conn == nil {
		panic("Connection has not been set up yet. Please call Setup() first")
	}

	return conn
}

func GenerateQueryAndModel() {
  g := gen.NewGenerator(gen.Config{
    OutPath: "./src/query",
		Mode: gen.WithoutContext|gen.WithDefaultQuery|gen.WithQueryInterface,
	})

	if conn == nil {
		panic("Connection has not been set up yet. Please call Setup() first")
	}

  g.UseDB(conn)
	g.ApplyBasic(g.GenerateAllTable()...)
  g.Execute()
}

func newConnection(user, password, host, port, dbname string) (*gorm.DB, error) {
	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
	connection, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{},
	)
	if (err != nil) {
		return nil, err
	}
	return connection, nil
}

