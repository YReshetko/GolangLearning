package main

import (
	"fmt"
	"mailclient/config"
	"mailclient/fetch"
	"os"
)

func main() {
	var config config.Configuration
	ok := config.Load()
	if ok != nil {
		fmt.Println(ok)
		os.Exit(1)
	}
	log(config)
	client := fetch.NewImapClient()
	client.Init(config)
	err := client.Login()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer client.Logout()

}

func log(object interface{}) {
	fmt.Printf("Logging object: %+v\n", object)
}
