/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"os"

	"github.com/mnsdojo/gotube/cmd"
)

func main() {
	if err:=cmd.Execute();err!=nil{
		fmt.Fprintln(os.Stderr,err)
		os.Exit(1)
	}
}

