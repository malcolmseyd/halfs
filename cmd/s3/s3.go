package main

import (
	"fmt"
	"os"

	"github.com/malcolmseyd/halfs/store"
)

const help = `Usage: s3 [COMMAND]

Commands:
	put <FILE>          stores FILE remotely and prints a reference name
	get <NAME> <FILE>   retrieves NAME and writes it to FILE`

func main() {
	if len(os.Args) < 2 {
		fmt.Println(help)
		os.Exit(1)
	}
	cmd := os.Args[1]
	args := os.Args[2:]

	if cmd == "put" {
		if len(args) < 1 {
			fmt.Println("Missing FILE argument")
			fmt.Println(help)
			os.Exit(1)
		}
		filename := args[0]
		blob, err := os.ReadFile(filename)
		if err != nil {
			fmt.Println("File can't be read:", err)
			os.Exit(1)
		}
		ref, err := store.Put(blob)
		if err != nil {
			fmt.Println("Failed to store file:", err)
			os.Exit(1)
		}
		fmt.Println(ref)
	} else if cmd == "get" {
		if len(args) < 1 {
			fmt.Println("Missing NAME argument")
			fmt.Println(help)
			os.Exit(1)
		}
		ref := args[0]
		if len(args) < 1 {
			fmt.Println("Missing FILE argument")
			fmt.Println(help)
			os.Exit(1)
		}
		filename := args[1]
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println("Failed to create file:", err)
			os.Exit(1)
		}
		defer file.Close()
		blob, err := store.Get(ref)
		if err != nil {
			fmt.Println("Failed to retrieve file:", err)
			os.Exit(1)
		}
		_, err = file.Write(blob)
		if err != nil {
			fmt.Println("Failed to write to file:", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully wrote %v bytes to %v\n", len(blob), filename)
	} else {
		fmt.Println("Unknown command:", cmd)
		os.Exit(1)
	}
}
