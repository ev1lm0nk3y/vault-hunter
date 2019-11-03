package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/vault/api"
)

const secretRoot = "platform-prod-comm_secrets/apps/access-controller/secrets"

func main() {
	c, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatalf("Cannot create client: %s", err)
	}
	client := c.Logical()

	var paths *api.Secret
	paths, err = client.Read(secretRoot)
	if err != nil {
		fmt.Println(paths)
		fmt.Println(err)
	}
	fmt.Println(paths.Data)
}
