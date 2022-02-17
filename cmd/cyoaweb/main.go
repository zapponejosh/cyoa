package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joshzappone/cyoa"
)

func main() {
	filename := flag.String("file", "gopher.json", "JSON file with CYOA story!")
	flag.Parse()
	file, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}

	story, err := cyoa.JsonStory(file)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", story)

}
