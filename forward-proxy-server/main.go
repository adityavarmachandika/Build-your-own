package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	// "net/url"
	"log"
)

var infologger, errorLogger *log.Logger

func baseHandler(w http.ResponseWriter, r *http.Request) {

	targetHostAdd := r.RequestURI

	infologger.Println("requset object", r)
	infologger.Println("r.requesturi : ", targetHostAdd, "r.remoteaddr : ", r.RemoteAddr)
	if targetHostAdd == "" {
		errorLogger.Println("No host address present")
		w.Write([]byte("there is no host name in the request"))
		return
	}

	outRequest, err := http.NewRequest(r.Method, targetHostAdd, r.Body)

	if err!=nil{
		errorLogger.Println("invalid request method : ", r.Method, " request Uri : ", r.RequestURI, "request body : ", r.Body)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Incorrect method or uri"))
		return
	}
	
	hopByHopHeaderMap := map[string]bool{"Connection": true,
		"Keep-Alive":          true,
		"Proxy-Authenticate":  true,
		"Proxy-Authorization": true,
		"Te":                  true,
		"Trailer":             true,
		"Transfer-Encoding":   true,
		"Upgrade":             true,
	}
	for key, values := range r.Header {

		if hopByHopHeaderMap[key]{
			continue
		}
		for _, value := range values {
			outRequest.Header.Add(key, value)
		}
	}
	outRequest.Header.Add("X-Forwarded-For",r.RemoteAddr)
	outRequest.Host = outRequest.URL.Host

	//fmt.Println(outRequest.Header)
	
	//creates a client before sending a requset.
	client := &http.Client{}

	// 6. Send the request out to the open web
	resp, err := client.Do(outRequest)
	if err != nil {
		http.Error(w, "Proxy destination unreachable: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close() 

	for key, values := range resp.Header {

		if hopByHopHeaderMap[key]{
			continue
		}
		for _, val := range values {
			w.Header().Add(key, val)
		}
	}

	w.WriteHeader(resp.StatusCode)
	//copies the response from the client. once done response is returned.
	_, err = io.Copy(w, resp.Body)


	if err != nil {
		fmt.Println("there is an error while sending the response", err)
	}

}
func main() {

	logFile, err := os.OpenFile("forwardProxyServer.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println("logfile not created. server is stopped")
		return
	}
	infologger = log.New(logFile, "INFO:", log.Lshortfile)
	errorLogger = log.New(logFile, "ERROR:", log.Lshortfile)
	port := ":8080"

	fmt.Println("...Starting Forward-Proxy-Server...\n Port : ", port)

	mux := http.NewServeMux()
	mux.HandleFunc("/", baseHandler)

	err = http.ListenAndServe(port, mux)
	if err != nil {
		errorLogger.Println("Error while starting a server : ", err)
	}

}
