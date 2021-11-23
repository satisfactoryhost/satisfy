package main

import "github.com/satisfactoryhost/satisfy/cmd"

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
