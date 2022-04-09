package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

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
	tpl := template.Must(template.New("").Parse(AltTmpl))

	h := cyoa.NewHandler(story, cyoa.AltPathFn(PathFn), cyoa.WithTemplate(tpl))
	h2 := cyoa.NewHandler(story)

	mux := http.NewServeMux()
	mux.Handle("/story/", h)
	mux.Handle("/", h2)
	fmt.Printf("Starting server on http://localhost:%d\n", *port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}
func PathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story/" || path == "/story" {
		path = "/story/intro"
	}

	return path[len("/story/"):]
}

var AltTmpl = `
<!DOCTYPE html>
<html lang="en">
  <head>
	<link rel="stylesheet" href="https://unpkg.com/purecss@2.0.6/build/pure-min.css" integrity="sha384-Uu6IeWbM+gzNVXJcM9XV3SohHtmWE+3VGi496jvgX1jyvDTXfdK+rfZc8C1Aehk5" crossorigin="anonymous">
    <meta charset="UTF-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Choose Your Own Adventure</title>
  </head>
  <body style="text-align: center; max-width: 900px; margin: 0 auto;">
	<h1 class="">{{.Title}} ALT</h1>
	<div class="">
    
    {{ range .Paragraphs}}
    <p class="">{{.}}</p>
    {{end}}

    <ul class="" style="margin: 0 auto;">
      {{range .Options}}
      <a class="pure-button" style="margin:5px;" href="/story/{{.Chapter}}">{{.Text}}</a>
      {{end}}
    </ul>
		</div>
  </body>
</html>`
