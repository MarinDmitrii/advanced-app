package main

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"os"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	file, err := os.OpenFile(os.Getenv("logfile"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile")
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.SetLevel(log.InfoLevel)
	//log.SetReportCaller(true) //добавляет номер строки
}

func main() {
	http.HandleFunc("/healthz", Healthz)
	http.HandleFunc("/user", createUser)

	makeHTTPserver()
}

func makeHTTPserver() {
	err := http.ListenAndServe(os.Getenv("localURL"), nil)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("It's alive!"))
		if err != nil {
			return
		}
	}
}

type user struct {
	Name  string `json:"Name"`
	Email string `json:"Email"`
}

func createUser(w http.ResponseWriter, r *http.Request) {
	headerContentType := r.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}
	var e user
	var unmarshalErr *json.UnmarshalTypeError

	//body, err := io.ReadAll(r.Body)
	//json.Unmarshal(body, &e)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&e)

	if err != nil {
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}
	userResponse(w, e, http.StatusOK)
	//log.Printf("%s", r.RemoteAddr)
	//log.Printf("%s", r.Header)

	reqDump, err1 := httputil.DumpRequest(r, true)
	if err1 != nil {
		log.Fatalf("%s", err1)
	}
	//log.Debug(string(reqDump))
	log.Print(string(reqDump))

	return
}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	log.Printf("%s", jsonResp)
}

func userResponse(w http.ResponseWriter, e user, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	jsonUser, err := json.Marshal(&e)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonUser)
	log.Printf("%s", jsonUser)
}
