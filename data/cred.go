package data

import "fmt"

type Credentials struct {
	Username string
	Password string
}

func Store(credMap map[string]Credentials, site string, username string, password string) error {
	if credMap == nil {
		return fmt.Errorf("credMap cannot be nil")
	}
	if site == "" || username == "" || password == "" {
		return fmt.Errorf("site, username, and password cannot be empty")
	}
	credMap[site] = Credentials{username, password}
	return nil
}
