package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joshzappone/cyoa"
)

func main() {
	port := flag.Int("port", 3000, "Port to start the applicaion")
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

	h := cyoa.NewHandler(story)

	fmt.Printf("Starting server on http://localhost:%d\n", *port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}
