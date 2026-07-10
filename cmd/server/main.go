package main

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"
	"slices"
	"strings"

	"golang.org/x/net/websocket"
)

//go:embed web/static
var staticFiles embed.FS

//go:embed web/templates
var templateFiles embed.FS

type NotesPageData struct {
	ID      string
	Content template.HTML
}

func notesPage(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile("./test-file.md")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := r.PathValue("id")
	tmpl, err := template.ParseFS(templateFiles, "web/templates/index.tmpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := NotesPageData{
		ID:      id,
		Content: template.HTML(MDToHTML(content)),
	}
	tmpl.Execute(w, data)
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func wsHandler(ws *websocket.Conn) {
	for {
		var msg string
		if err := websocket.Message.Receive(ws, &msg); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error: %v\n", err)
		}
		words := strings.Split(msg, " ")
		slices.Reverse(words)
		revMsg := strings.Join(words, " ")
		websocket.Message.Send(ws, revMsg)
	}
}

func main() {
	staticFS, err := fs.Sub(staticFiles, "web/static")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/notes/{id}", notesPage)
	http.Handle("/ws", websocket.Handler(wsHandler))
	http.HandleFunc("/health", health)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServerFS(staticFS)))

	fmt.Println("Listening on port :8080")
	http.ListenAndServe(":8080", nil)
}
