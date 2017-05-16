package main

import (
	"fmt"
	"os"
	"flag"
	
	"github.com/umayr/lunch"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("error occurred")
			os.Exit(1)
		}
	}()

	bare := flag.Bool("bare", false, "runs the bare minimum functionality, no authentication in any case")
	flag.Parse()

	l, err := lunch.Today(!*bare)
	if err != nil {
		fmt.Println(fmt.Sprintf("error occurred: %s", err.Error()))
		os.Exit(1)
	}

	fmt.Println(l.Name)
}
