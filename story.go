package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

func JsonStory(r io.Reader) (Story, error) {
	var story Story
	err := json.NewDecoder(r).Decode(&story)
	if err != nil {
		return nil, err
	}
	return story, nil
}

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

type HandlerOpt func(h *handler)

func WithTemplate(t *template.Template) HandlerOpt {
	return func(h *handler) {
		h.t = t
	}
}

func AltPathFn(fn func(*http.Request) string) HandlerOpt {
	return func(h *handler) { h.PathFn = fn }
}

type handler struct {
	s      Story
	t      *template.Template
	PathFn func(r *http.Request) string
}

func NewHandler(s Story, opts ...HandlerOpt) http.Handler {
	h := handler{s, tpl, defaultPathFn}

	for _, opt := range opts {
		opt(&h)
	}

	return h
}

func defaultPathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	path = path[1:]
	if path == "" {
		path = "intro"
	}
	return path
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.PathFn(r)

	if chapter, ok := h.s[path]; ok {
		err := h.t.Execute(w, chapter)

		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Oh no! You've encountered an error!", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "No chapter found :/", http.StatusNotFound)
}

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

var tpl *template.Template

var defaultHandlerTmpl = `
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
	<h1 class="">{{.Title}}</h1>
	<div class="">
    
    {{ range .Paragraphs}}
    <p class="">{{.}}</p>
    {{end}}

    <ul class="" style="margin: 0 auto;">
      {{range .Options}}
      <a class="pure-button" style="margin:5px;" href="/{{.Chapter}}">{{.Text}}</a>
      {{end}}
    </ul>
		</div>
  </body>
</html>`
