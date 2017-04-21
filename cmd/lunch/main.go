package main

import (
	"fmt"
	"os"

	"github.com/umayr/lunch"
)

func main() {
	l, err := lunch.Today()
	if err != nil {
		fmt.Println(fmt.Sprintf("error occurred: %s", err.Error()))
		os.Exit(1)
	}

	fmt.Println(l.Name)
}
