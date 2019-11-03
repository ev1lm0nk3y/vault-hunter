package main

import (
	"filepath"
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/hashicorp/vault/api"
)

var (
	secretRoot string
	value      string
)

func vaultClient() (*api.Logical, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}
	return client.Logical(), nil
}

func (c *api.Logical) listSecretPath(path string) error {
	wg := *sync.WaitGroup()
	paths, err := c.List(path)
	if err != nil {
		return err
	}

	if _, ok := paths.Data["keys"]; ok {
		for p := range paths.Data["keys"] {
			newPath := filepath.Join(path, p)
			if strings.Contains(p, "/") {
				wg.Add(1)
				go listSecretPath(newPath)
			} else {
				wg.Add(1)
				go c.valueInSecret(newPath, wg)
			}
		}
	}
	wg.Wait()
	return nil
}

func (c *api.Logical) valueInSecret(path string, wg *sync.WaitGroup) {
	defer wg.Done()
	s, err := c.Read(path)
	if err != nil {
		return
	}

	for _, v := range s.Data {
		if v == value {
			fmt.Println(path)
		}
	}
}

func main() {
	flag.StringVar(&secretRoot, "s", "", "Vault secret root to start searching")
	flag.StringVar(&value, "v", "", "Value to search for in vault")

	client, err := vaultClient()
	if err != nil {
		log.Fatalf(err)
	}

	err = listSecretPath(secretRoot)

	fmt.Println("vim-go")
}
