package main

//this contains the interface (for now it will be connected to terminal, later C++ Qt will connect to it)

import (
	"fmt"

	"github.com/rihts-4/pasword-mango/data"
)

func printAllCred(credMap map[string]data.Credentials) {
	for site, credentials := range credMap {
		fmt.Printf("Site: %s, Username: %s, Password: %s\n", site, credentials.Username, credentials.Password)
	}
}

func main() {
	credMap := make(map[string]data.Credentials)
	data.Store(credMap, "google", "admin", "password")
	printAllCred(credMap)
	data.Store(credMap, "fabook", "nouse", "p125")
	printAllCred(credMap)
	fmt.Println("Hello, Go Projec!")
}
