package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	// "net/url"
	"log"
)

var infologger, errorLogger *log.Logger
var forbiddenHosts, bannedWords map[string]bool

func checkForBannedWords(body []byte) bool {
	newbody := strings.ToLower(string(body))

	for word := range bannedWords {
		if strings.Contains(newbody, word) {
			infologger.Println("banned word found in the response : ", word)
			return true
		}
	}
	return false
}

func loadForbids(fileName string) (map[string]bool, error) {

	wordMap := make(map[string]bool)

	file, err := os.Open(fileName)
	if err != nil {
		errorLogger.Println("unable to load the file. filename : ", fileName, " err : ", err)
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		line := strings.TrimSpace(scanner.Text())

		if line != "" && !strings.HasPrefix(line, "//") {
			wordMap[line] = true
		}
	}
 
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file contents: %v\n", err)
		return nil, err
	}

	return wordMap, nil
}

func baseHandler(w http.ResponseWriter, r *http.Request) {

	//infologger.Printf("%+v\n", r)
	if r.Method == http.MethodConnect {


		hostName:=strings.TrimPrefix(r.URL.Hostname(),"www.")
		if forbiddenHosts[hostName] {
			infologger.Println("A forbidden host requested")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("forbidden host\n"))
			return
		}

		//establish a tcp conncection to the client requested address
		destConn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)

		if err != nil {
			errorLogger.Println("there is a error while creating the making a call to the client requested server", err)
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		defer destConn.Close()

		hijack, ok := w.(http.Hijacker)

		if !ok {
			errorLogger.Println("web server not able to create a hijact")
			http.Error(w, "web server not able to create hijack interface", http.StatusInternalServerError)
			return
		}

		clientConn, _, err := hijack.Hijack()

		if err != nil {
			errorLogger.Println("there is an error while creating client conn", err)
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		defer clientConn.Close()

		_, err = clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

		if err != nil {
			errorLogger.Println("the client connection returned an error", err)
			return
		}

		errChan := make(chan error, 2)
		go func() {
			_, err := io.Copy(destConn, clientConn)
			errChan <- err
		}()
		go func() {
			_, err := io.Copy(clientConn, destConn)
			errChan <- err
		}()

		// Wait for one side of the connection to terminate or fail
		<-errChan

	} else {

		infologger.Println("-Request from client : ", r.RemoteAddr, " request url : ", r.URL, " request method : ", r.Method)

		targetHostAdd := r.RequestURI

		//infologger.Printf("%+v\n", r)
		//infologger.Println("r.requesturi : ", targetHostAdd, "r.remoteaddr : ", r.RemoteAddr)
		if targetHostAdd == "" {
			errorLogger.Println("No host address present")
			w.Write([]byte("there is no host name in the request"))
			return
		}

		//return if the request is to a forbidden host
		if forbiddenHosts[r.Host] {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("forbidden host\n"))
			return
		}

		outRequest, err := http.NewRequest(r.Method, targetHostAdd, r.Body)

		if err != nil {
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

			if hopByHopHeaderMap[key] {
				continue
			}
			for _, value := range values {
				outRequest.Header.Add(key, value)
			}
		}
		outRequest.Header.Add("X-Forwarded-For", r.RemoteAddr)
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

		var tempBody []byte
		tempBody, err = io.ReadAll(resp.Body)

		if err != nil {
			fmt.Println("there is an error while sending the response", err)
		}

		containsBannedWords := checkForBannedWords(tempBody)
		if containsBannedWords {
			http.Error(w, "response not safe", http.StatusBadGateway)
			return
		}
		for key, values := range resp.Header {

			if hopByHopHeaderMap[key] {
				continue
			}
			for _, val := range values {
				w.Header().Add(key, val)
			}
		}

		w.WriteHeader(resp.StatusCode)
		//copies the response from the client. once done response is returned.
		_, err = w.Write(tempBody)

		if err != nil {
			fmt.Println("there is an error while sending the response", err)
		}

		infologger.Println(r.RemoteAddr, " ", resp.StatusCode, " OK")

	}

}
func main() {
	var err error
	forbiddenHosts, err = loadForbids("forbidden-hosts.txt")
	if err != nil || forbiddenHosts == nil {
		fmt.Println("error while fetching forbidden hosts map", err)
		return
	}

	bannedWords, err = loadForbids("banned-words.txt")
	if err != nil || bannedWords == nil {
		fmt.Println("error while fetching banned words map")
		return
	}

	logFile, err := os.OpenFile("forwardProxyServer.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println("logfile not created. server is stopped")
		return
	}
	infologger = log.New(logFile, "INFO:", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(logFile, "ERROR:", log.Ldate|log.Ltime|log.Lshortfile)
	port := ":8080"

	fmt.Println("...Starting Forward-Proxy-Server...\n Port : ", port)

	err = http.ListenAndServe(port, http.HandlerFunc(baseHandler))
	if err != nil {
		errorLogger.Println("Error while starting a server : ", err)
	}

}
