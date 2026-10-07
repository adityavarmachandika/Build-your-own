package main

import (
	"fmt"
	"io"
	"net/http"
	// "net/url"
)


func baseHandler(w http.ResponseWriter, r *http.Request){

	fmt.Println(*r)

	targetHostAdd:= r.RequestURI;

	if targetHostAdd == ""{
		w.Write([]byte("there is no host name in the reqest"))
		return
	}

	//using net/url package to parse the hosturl
	// parsedURL,err:= url.Parse(targetHostAdd)

	// if err != nil{
	// 	panic(err)
	// }

	//this is the parsed url
	// fmt.Printf("%v\n",parsedURL.Scheme)
	// fmt.Printf("%v\n",parsedURL.Host)
	// fmt.Printf("%v\n",parsedURL.Fragment)
	// fmt.Printf("%v\n",*parsedURL)

	//need to do a url validation here.

	outRequest,err:= http.NewRequest(r.Method,targetHostAdd,r.Body)

	fmt.Println(targetHostAdd)

	outRequest.Header=r.Header.Clone()

	outRequest.Host=targetHostAdd

	//creates a client before sending a requset.
	client := &http.Client{}


	// 6. Send the request out to the open web
	resp, err := client.Do(outRequest)
	if err != nil {
		http.Error(w, "Proxy destination unreachable: "+err.Error(), http.StatusBadGateway)
		return
	}

	for key,values :=range resp.Header{
		for _,val :=range values{
			w.Header().Add(key,val)
		}
	}
	_,err=io.Copy(w,resp.Body)

	if err != nil{
		fmt.Println("there is an error while sending the response", err)
	}

}
func main() {
	fmt.Println("Hello world")

	port:=":8080"

	mux := http.NewServeMux()
	fmt.Printf("%v\n",mux)	

	mux.HandleFunc("/",baseHandler)
	fmt.Println("this is the port number", port)

	fmt.Println(http.ListenAndServe(port,mux))

}