package main 

import (
	"fmt"
	"net/http"
)
func bikeHandler(w http.ResponseWriter, r *http.Request){
	title := r.URL.Path[len("/view/"):]
	b, _ := bikeLoad(title)
	fmt.Printf("%T,%v \n",b,b)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s<p>%v</p></div>", b.Name, b.Brand, b.Year)

}

func main() {
	http.HandleFunc("/view/",bikeHandler)
	http.ListenAndServe(":8080", nil)
	fmt.Printf("Listening To :8080")
}