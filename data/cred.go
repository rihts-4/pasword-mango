package data

import (
	"fmt"
	"sync"
)

type Credentials struct {
	Username string
	Password string
}

var mu sync.RWMutex
var credMap map[string]Credentials = make(map[string]Credentials) //consider race conditions (later iterations)
var siteSearch *Trie                                              //consider race conditions (later iterations)

func InitCredMap() {
	credMap = make(map[string]Credentials)
	siteSearch = NewTrie()
} //consider concurrency issues

func Store(site string, username string, password string) {
	mu.Lock()
	defer mu.Unlock()

	if credMap == nil {
		credMap = make(map[string]Credentials)
	}
	credMap[site] = Credentials{username, password}
	siteSearch.Insert(site)
}

func Update(site string, username string, password string) {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := Retrieve(site); !exists {
		fmt.Println("No such site to update.")
		return
	}

	credMap[site] = Credentials{username, password}
	fmt.Println("Credentials updated successfully!")
}

func Show() {
	mu.RLock()
	defer mu.RUnlock()

	for site, credentials := range credMap {
		fmt.Printf("Site: %s, Username: %s, Password: %s\n", site, credentials.Username, credentials.Password)
	}
}

func Retrieve(site string) (Credentials, bool) {
	mu.RLock()
	defer mu.RUnlock()

	if credMap == nil {
		fmt.Println("Map not initialized.")
	}

	if siteSearch.SearchWord(site) {
		return credMap[site], true
	}
	return Credentials{}, false
}

func Delete(site string) bool {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := Retrieve(site); exists {
		delete(credMap, site)
		siteSearch.Delete(site) // Also delete from the trie
		fmt.Println("Credentials deleted successfully!")
		return true
	}
	return false
}
