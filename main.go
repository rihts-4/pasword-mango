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
	credMap := make(map[string]data.Credentials)
	fmt.Println("Interface Starts: Running....")
	data.Store(credMap, "google", "admin", "passsword")
	data.Store(credMap, "fabook", "nouse", "p125")
	data.Show(credMap)
	data.Update(credMap, "google", "admin_updated", "newpassword")
	data.Show(credMap)
	fmt.Println("Interface Ends: Exiting....")
}
