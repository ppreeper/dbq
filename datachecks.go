package main

import (
	"fmt"
	"os"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func fatalErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(10)
	}
}
