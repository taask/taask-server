package main

import "log"

func main() {
	root := rootCmd()
	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
