package main

import (
	"fmt"
	"net/http"
)

func baseHandler(w http.ResponseWriter, r *http.Request){

	w.Write([]byte("this response is from port simple-server\n"))
}

func main (){
	fmt.Println("hello 8989")

	port:= ":8989"

	mux:=http.NewServeMux()

	mux.HandleFunc("/", baseHandler)
	http.ListenAndServe(port,mux)

}