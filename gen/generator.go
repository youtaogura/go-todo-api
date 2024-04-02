package main

import (
	"go_todo/src/database"
)

func main() {
	database.Setup()
	database.GenerateQueryAndModel()
}
