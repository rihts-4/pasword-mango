package data

import "fmt"

type Credentials struct {
	Username string
	Password string
}

var credMap map[string]Credentials = make(map[string]Credentials)
var siteSearch *Trie

func InitCredMap() {
	credMap = make(map[string]Credentials)
	siteSearch = NewTrie()
}

func Store(site string, username string, password string) {
	if credMap == nil {
		credMap = make(map[string]Credentials)
	}
	credMap[site] = Credentials{username, password}
	siteSearch.Insert(site)
}

func Update(site string, username string, password string) {
	if _, exists := Retrieve(site); !exists {
		fmt.Println("No such site to update.")
		return
	}

	credMap[site] = Credentials{username, password}
	fmt.Println("Credentials updated successfully!")
}

func Show() {
	for site, credentials := range credMap {
		fmt.Printf("Site: %s, Username: %s, Password: %s\n", site, credentials.Username, credentials.Password)
	}
}

func Retrieve(site string) (Credentials, bool) {
	if credMap == nil {
		fmt.Println("Map not initialized.")
	}

	if siteSearch.SearchWord(site) {
		return credMap[site], true
	}
	return Credentials{}, false
}

func Delete(site string) bool {
	if _, exists := Retrieve(site); exists {
		delete(credMap, site)
		siteSearch.Delete(site) // Also delete from the trie
		fmt.Println("Credentials deleted successfully!")
		return true
	}
	return false
}
