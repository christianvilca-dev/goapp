package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//Page sirve para la estructura de las paginas
type Page struct {
	Title string
	Body  []byte // en vez de string para que sea mas eficiente
}

func (p *Page) save() error {
	filename := "./data/" + p.Title + ".txt"
	err := ioutil.WriteFile(filename, p.Body, 0600) // el 3er parametro son los permisos (solo puede leer y escribir el dueño del archivo)
	return err
}

func loadPage(title string) (*Page, error) {
	filename := "./data/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	page := &Page{Title: title, Body: body}
	return page, err
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):] // captura lo que esta despues del "/view/"
	p, err := loadPage(title)
	if err != nil {
		log.Fatal("Error al cargar pagina")
		//fmt.Fprintf(w, "<div>%s</div>", err)
	}
	// Fprintf con la interfaz write escribir en el
	fmt.Fprintf(w, "<h1>%s<h1/><div>%s</div>", p.Title, p.Body) //porque el p.body es byte y ya no hacemos una conversion
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	page, err := loadPage(title)
	if err != nil {
		page = &Page{Title: title} // Puntero a una estructura page
	}
	fmt.Fprintf(w, `
		<html>
		<head>
			<title>%s</title>
		</head>
		<body>
			<h1>%s</h1>
			<form method="POST" action="/save/%s">
				<textarea name="body">%s</textarea>
				<button>Guardar</button
			</form>
		</body>
		</html>	
	`, page.Title, page.Title, page.Title, page.Body)
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Bienvenidos</h1>")
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body") // Capturamos los elementos del formulario
	page := &Page{Title: title, Body: []byte(body)}
	page.save()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	/* page := &Page{Title: "primer", Body: []byte("Nuestra primer pagina")}
	page.save() */
	/* page := loadPage("primer")
	fmt.Println(page.Title, string(page.Body)) */

	// Se crea un servidor
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	fmt.Println("Ingrese a http://localhost:8080/view/")
	http.ListenAndServe(":8080", nil)
}
