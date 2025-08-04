package database

import (
	"fmt"
	"strings"
)

func InitDatatabase() {
	dbAddr := "localhost:5432"
	hostPort := strings.Split(dbAddr, ":")
	dbUser := "test_user"
	dbPassword := "test_password"
	dbName := "test_db"
	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=disable", hostPort[0], dbUser, dbPassword, hostPort[1], dbName)

	fmt.Println("Connecting to database with DSN:", dsn)

}
