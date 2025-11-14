package data

import "fmt"

type Credentials struct {
	Username string
	Password string
}

func Store(credMap map[string]Credentials, site string, username string, password string) {
	credMap[site] = Credentials{username, password}
	fmt.Println("Credentials stored successfully!")
}
