package main

//this contains the interface (for now it will be connected to terminal, later C++ Qt will connect to it)
/*
Current Assumptions:
- One User (me) ## SO NO USER AUTH
- Small DB
*/

import (
	"fmt"

	"github.com/rihts-4/pasword-mango/data"
)

func main() {
	fmt.Println("Interface Starts: Running....")
	data.InitCredMap()
	//Data Population
	data.Store("google", "admin", "password")
	data.Store("facebook", "nouse", "p125")
	data.Store("example", "user", "password123")
	data.Store("twitter", "john_doe", "Tw1tt3r2024!")
	data.Store("github", "developer", "c0d3Master#99")
	data.Store("linkedin", "professional", "Career$ecure1")
	data.Store("instagram", "photo_user", "Insta@Pass567")
	data.Store("microsoft", "office_admin", "M$0ffice2024")
	data.Store("amazon", "shopper123", "Amaz0n!Prime")
	data.Store("netflix", "movie_buff", "Netfl1x&Chill")
	data.Store("spotify", "music_lover", "Spot1fy!Beats")
	data.Store("reddit", "redditor42", "R3dd!tK@rma")
	data.Store("dropbox", "cloud_user", "Dr0pb0x$ync!")
	//Data Population Ends

	credentials, found := data.Retrieve("github")
	if found {
		fmt.Printf("Credentials for github - Username: %s, Password: %s\n", credentials.Username, credentials.Password)
	} else {
		fmt.Println("No credentials found for github")
	}

	data.Update("github", "dev_guru", "N3wP@ssw0rd!")
	credentials, found = data.Retrieve("github")
	if found {
		fmt.Printf("Updated Credentials for github - Username: %s, Password: %s\n", credentials.Username, credentials.Password)
	} else {
		fmt.Println("No credentials found for github")
	}

	deleted := data.Delete("twitter")
	if deleted {
		fmt.Println("Twitter credentials deleted.")
	} else {
		fmt.Println("No Twitter credentials to delete.")
	}

	data.Show()

	fmt.Println("Interface Ends: Exiting....")
}
