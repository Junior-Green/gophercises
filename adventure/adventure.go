package adventure

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

const templatePath string = "/Users/juniorgreen/Documents/go_exercises/adventure/templates/page.html"
const staticPath string = "adventure/static"

type StoryArcName string

type StoryArcOptions []StoryArcOption

type StoryArcOption struct {
	Text string       `json:"text"`
	Arc  StoryArcName `json:"arc"`
}

type StoryArc struct {
	Title   string          `json:"title"`
	Story   []string        `json:"story"`
	Options StoryArcOptions `json:"options"`
}

type Stories map[StoryArcName]StoryArc

func (stories Stories) String() string {
	var str string
	for k, story := range stories {
		str += fmt.Sprintf("%s: %+v\n\n", k, story)
	}
	return str
}

func Init() {
	var (
		filename string
		err      error
	)

	//Parse command line flags
	flag.StringVar(&filename, "f", "adventure/goper.json", "References file that contains story data.")
	flag.Parse()

	//Read and parse JSON
	bytes, err := os.ReadFile(filename)
	handleError(err)
	stories, err := parseJSON(bytes)
	handleError(err)

	//Parse HTML template
	template, err := template.New("page.html").ParseFiles(templatePath)
	handleError(err)

	//Generate dynamic request handler
	handler := createMapHandler(stories, template)

	//Attach handlers
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))
	http.HandleFunc("/", handler)

	//Start server
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", nil)

}

func createMapHandler(stories Stories, template *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		arcName := r.URL.Path[1:]

		if storyArc, ok := stories[StoryArcName(arcName)]; ok {
			template.Execute(w, storyArc)
			return
		}

		http.Redirect(w, r, "/intro", http.StatusSeeOther)
	}
}

func parseJSON(jsn []byte) (Stories, error) {
	var stories Stories

	err := json.Unmarshal(jsn, &stories)

	if err != nil {
		return nil, err
	}

	return stories, nil
}

func handleError(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(-1)
	}
}
