package main

import (
	"fmt"
	"github.com/tgrangeray/nirvana-goclient/nirvana"
	"os"
)

func main() {

	username := os.Getenv("NIRVANA_USERNAME")
	password := os.Getenv("NIRVANA_PASSWORD")
	if len(username) == 0 || len(password) == 0 {
		fmt.Print("Set environment variables NIRVANA_USERNAME and NIRVANA_PASSWORD")
		return
	}

	nirvanaCli, _ := nirvana.NewNirvanaClient(nil)
	defer nirvanaCli.Close()

	err := nirvanaCli.Authenticate(username, password)
	if err != nil {
		fmt.Print(err)
		return
	}

	_, err = nirvanaCli.RetrieveSince(0)
	if err != nil {
		fmt.Print(err)
	}

}
