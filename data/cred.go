package data

import "fmt"

type Credentials struct {
	Username string
	Password string
}

func Store(credMap map[string]Credentials, site string, username string, password string) {
	if credMap == nil {
		credMap = make(map[string]Credentials)
	}
	credMap[site] = Credentials{username, password}
	fmt.Println("Credentials stored successfully!")
}

func Update(credMap map[string]Credentials, site string, username string, password string) {
	if _, exists := credMap[site]; !exists {
		fmt.Println("No such site to update.")
		return
	}

	credMap[site] = Credentials{username, password}
	fmt.Println("Credentials updated successfully!")
}

func Show(credMap map[string]Credentials) {
	for site, credentials := range credMap {
		fmt.Printf("Site: %s, Username: %s, Password: %s\n", site, credentials.Username, credentials.Password)
	}
}

func Retrieve(credMap map[string]Credentials, site string) (Credentials, bool) {
	if credMap == nil {
		fmt.Println("No such map")
	}

	cred, exists := credMap[site]
	if !exists {
		fmt.Println("No such site")
	}

	fmt.Println("Credentials retrieved successfully!")

	return cred, exists
}

func Delete(credMap map[string]Credentials, site string) bool {
	if _, exists := credMap[site]; exists {
		delete(credMap, site)
		fmt.Println("Credentials deleted successfully!")
		return true
	}
	return false
}
