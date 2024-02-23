package main

import (
	"fmt"
	"net/http"
)

var portNumber = ":8080"

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/about", About)

	fmt.Println("Run Web Application at localhost" + portNumber)
	_ = http.ListenAndServe(portNumber, nil)
}

func Home(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "This is Home Page!")
}

func About(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "This is About Page!")
}

func addValue(x, y int) (int, error) {
	return x + y, nil
}
