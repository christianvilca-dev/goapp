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
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		log.Fatal("Error al cargar pagina")
		//fmt.Fprintf(w, "<div>%s</div>", err)
	}
	// Fprintf con la interfaz write escribir en el
	fmt.Fprintf(w, "<h1>%s<h1/><div>%s</div>", p.Title, p.Body) //porque el p.body es byte y ya no hacemos una conversion
}

func main() {
	/* page := &Page{Title: "primer", Body: []byte("Nuestra primer pagina")}
	page.save() */
	/* page := loadPage("primer")
	fmt.Println(page.Title, string(page.Body)) */

	// Se crea un servidor
	http.HandleFunc("/view/", viewHandler)
	http.ListenAndServe(":8080", nil)
}
