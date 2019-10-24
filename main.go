package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
)

func main() {
	initRoute()
	startServer()
}

func startServer() {
	port := ":9000"
	fmt.Printf("Starting server at %s \n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func initRoute() {
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/", index)

	http.HandleFunc("/template", withTemplateBasic)

	http.HandleFunc("/template/partials-hello", withTemplatePartialsHello)

	http.HandleFunc("/template/partials-world", withTemplatePartialsWorld)

	http.HandleFunc("/template/partials-index", withTemplatePartialsButUseParseFile)

	http.HandleFunc("/template/action-and-variable", withTemplateActionAndVariable)
}

func loadTemplate(fileName string) *template.Template {
	var filePath = path.Join(
		"views/basic",
		fmt.Sprintf("%s.html", fileName))
	var tmpl, err = template.ParseFiles(filePath)
	if err != nil {
		panic(err.Error())
		return nil
	}
	return tmpl
}

// Metode parsing menggunakan template.ParseGlob()
// memiliki kekurangan yaitu sangat tergantung terhadap pattern path yang digunakan.
func loadTemplateWithPartials() *template.Template {
	var tmpl, err = template.ParseGlob("views/partials/*")
	if err != nil {
		panic(err.Error())
		return nil
	}
	return tmpl
}

type M map[string]interface{}

type Info struct {
	Affiliation string
	Address     string

}

type Person struct {
	Name    string
	Gender  string
	Hobbies []string
	Info    Info
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, GO"))
}

func withTemplateBasic(w http.ResponseWriter, r *http.Request) {
	tmpl := loadTemplate("index")
	var data = M{
		"title": "Just test",
		"message": "Whops that's amazing",
	}
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func withTemplatePartialsHello(w http.ResponseWriter, r *http.Request) {
	var data = M{
		"title": "Test with partials",
		"message": "Oh yea partials from hello",
	}
	tmpl := loadTemplateWithPartials()
	err := tmpl.ExecuteTemplate(w, "hello", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func withTemplatePartialsWorld(w http.ResponseWriter, r *http.Request) {
	var data = M{
		"title": "Test with partials",
		"message": "Oh yea partials from World",
	}
	tmpl := loadTemplateWithPartials()
	err := tmpl.ExecuteTemplate(w, "world", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func withTemplatePartialsButUseParseFile(w http.ResponseWriter, r *http.Request) {
	var data = M{
		"title": "Test with partials",
		"message": "Oh yea partials with parse file",
	}

	//Fungsi ini selain bisa digunakan untuk parsing satu buah file saja
	var tmpl = template.Must(template.ParseFiles(
		"views/partials/index.html",
		"views/partials/_header.html",
		"views/partials/_footer.html",
	))

	var err = tmpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func withTemplateActionAndVariable(w http.ResponseWriter, r *http.Request) {
	var person = Person{
		Name:    "Bruce Wayne",
		Gender:  "male",
		Hobbies: []string{"Reading Books", "Traveling", "Buying things"},
		Info:    Info{"Wayne Enterprises", "Gotham City"},
	}
	var tmpl = template.Must(template.ParseFiles("views/basic/action-and-var.html"))
	if err := tmpl.Execute(w, person); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (t Info) GetAffiliationDetailInfo() string {
	return "have 31 divisions"
}