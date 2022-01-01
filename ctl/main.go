package main

import "log"

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
