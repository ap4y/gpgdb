package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/howeyc/gopass"

	"github.com/ap4y/gpgdb/lib"
)

var client *Client

const usage = `Usage: %s [options] command [key] [value]

Commands:
  put     Put *key* with *value*
  keys    List all available keys
  get     Get value of the *key*
  delete  Delete value of the *key*

Options:
`

func main() {
	flag.Usage = func() {
		fmt.Printf(usage, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	usr, _ := user.Current()
	keyring := flag.String("keyring", usr.HomeDir+"/.gpgdb/secring.gpg", "path to the secret keyring")
	host := flag.String("host", "http://localhost:8080", "daemon host")
	flag.Parse()

	entity, err := lib.NewEntity(*keyring, func(entity *lib.Entity) ([]byte, error) {
		fmt.Print("Please enter keyring passphrase: ")
		passphrase := gopass.GetPasswd()
		return passphrase, nil
	})
	if err != nil {
		fmt.Printf("Unable to decrypt the keyring: %s", err)
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Printf("Missing command argument")
		return
	}

	client = NewClient(*host, entity)

	switch args[0] {
	case "put":
		put(args[1:])
	case "keys":
		keys()
	case "get":
		get(args[1:])
	case "delete":
		delete(args[1:])
	}
}

func put(args []string) {
	if len(args) != 2 {
		fmt.Printf("Missing key or value argument")
		return
	}

	if err := client.Put(args[0], args[1]); err != nil {
		fmt.Printf("Unable to put value: %s", err)
	} else {
		log.Printf("Value for key %s added", args[0])
	}
}

func keys() {
	if keys, err := client.Keys(); err != nil {
		fmt.Printf("Unable to get keys: %s", err)
	} else {
		log.Printf("Keys: %s", keys)
	}
}

func get(args []string) {
	if len(args) != 1 {
		fmt.Printf("Missing key argument")
		return
	}

	if value, err := client.Get(args[0]); err != nil {
		fmt.Printf("Unable to get value: %s", err)
	} else {
		fmt.Printf("Value: %s", value)
	}
}

func delete(args []string) {
	if len(args) != 1 {
		fmt.Printf("Missing key argument")
		return
	}

	if err := client.Delete(args[0]); err != nil {
		fmt.Printf("Unable to delete value: %s", err)
	} else {
		log.Printf("Value for key %s removed", args[0])
	}
}
