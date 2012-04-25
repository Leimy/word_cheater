package main

import (
	"fmt"
	ana "github.com/Leimy/anagrammer"
	ss "github.com/Leimy/sortstring"
	"log"
	"net/http"
//	"old/template"
	"html/template"
	"sort"
	"strings"
)

var ts = make(map[string]*template.Template)

var anagrams map[string][]string

// func init() {
// 	anagrams = ana.AnagramsFromFile("/usr/share/dict/words")
// 	for _, tmpl := range []string{"input", "results"} {
// 		ts[tmpl] = template.MustParseFile(tmpl+".html", nil)
// 	}
// }

func init () {
	anagrams = ana.AnagramsFromFile("/usr/share/dict/words")
	for _, tmpl := range []string{"input", "results"} {
		ts[tmpl] = template.Must(template.ParseFiles(tmpl+".html"))
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *templateParams) {
	err := ts[tmpl].Execute(w, p)
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
	}
}

// Could have just put the HTML in this function, but I like templating for later
// changes should I want them
func goHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "input", nil)
}

func getAllPerms(in ss.SortString) []string {
	sort.Sort(in)

	results := make([]string, 0)
	results = append(results, in.String())
	for in.NextPermutation(0, in.Len()) == true {
		results = append(results, in.String())

	}

	return results
}

// maps have unique keys, so lets just use it as a set
func getUniques(perms []string, size int) map[string]int {
	results := make(map[string]int)
	for _, cur := range perms {
		cur = cur[0:size] // slice it
		sorts := ss.NewSortString(cur)
		sort.Sort(sorts)            // sort it
		results[sorts.String()] = 0 // store it
	}

	return results
}

func runHandler(w http.ResponseWriter, r *http.Request) {
	in := ss.NewSortString(r.FormValue("input"))
	allperms := getAllPerms(in)
	size := in.Len()
	output := "Results are: \n"
	for n := size; n >= 3; n-- {
		output += fmt.Sprintf("\n\nFor size: %d\n", n)
		uniques := getUniques(allperms, n)
		for k, _ := range uniques {
			output += fmt.Sprintf("%s", strings.Join(anagrams[k], "\n"))
		}
	}
	output += fmt.Sprintf("\n")
	renderTemplate(w, "results", &templateParams{output, in.String()})
}

type templateParams struct {
	Results    string
	User_input string
}

func main() {
	http.HandleFunc("/run", runHandler)
	http.HandleFunc("/go", goHandler)

	log.Print("Ready...")

	http.ListenAndServe(":8080", nil)
}
