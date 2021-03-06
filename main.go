package main

import (
	"brightlocal/check"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

type Data struct {
	Method  string `json:"method,omitempty"`
	Key     string `json:"key,omitempty"`
	Value   string `json:"value,omitempty"`
	Error   string `json:"error,omitempty"`
	Result  string `json:"result,omitempty"`
	storage map[string]string
}

var data Data

func init() {
	data = Data{storage: make(map[string]string)}
}

func main() {
	http.HandleFunc("/req", req)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func req(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.Write([]byte(`{"error":"something went wrong"}`))
		return
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(err)
		w.Write([]byte(`{"error":"invalid json"}`))
		return
	}

	switch data.Method {
	case "GET":
		data.get(w)
	case "SET":
		data.set(w)
	case "DELETE":
		data.remove(w)
	case "EXISTS":
		data.exists(w)
	default:
		wrongMethod(w)
	}
}

func (d *Data) get(w http.ResponseWriter) {
	result := Data{
		Method: "GET",
		Key:    d.Key,
	}

	if check.KeyExists(d.Key, d.storage) {
		result.Value = d.storage[d.Key]
	} else {
		result.Error = "not found"
	}

	render(result, w)
}

func (d *Data) set(w http.ResponseWriter) {
	var mu sync.Mutex

	result := Data{
		Method: "SET",
		Value:  d.Value,
		Key:    d.Key,
	}

	if check.Value(d.Value) {
		result.Error = "value too long"
	} else if check.Key(d.Key) {
		result.Error = "key too long or empty"
	} else if check.StorageSize(d.storage) {
		result.Error = "storage is full. maximum size 1024 keys"
	} else {
		mu.Lock()
		d.storage[d.Key] = d.Value
		mu.Unlock()
	}

	render(result, w)
}

func (d *Data) exists(w http.ResponseWriter) {
	result := Data{
		Method: "EXISTS",
		Key:    d.Key,
	}

	if check.KeyExists(d.Key, d.storage) {
		result.Result = "exists"
	} else {
		result.Result = "not exists"
	}

	render(result, w)
}

func (d *Data) remove(w http.ResponseWriter) {
	result := Data{
		Method: "DELETE",
		Key:    d.Key,
	}

	if check.KeyExists(d.Key, d.storage) {
		delete(d.storage, d.Key)

		result.Result = "success"
	} else {
		result.Error = "not found"
	}

	render(result, w)
}

func wrongMethod(w http.ResponseWriter) {
	methods := []string{"GET", "SET", "DELETE", "EXISTS"}

	result := Data{
		Error: "method not allowed. allowed methods is: " + strings.Join(methods, ","),
	}

	render(result, w)
}

// render JSON response
func render(result Data, w http.ResponseWriter) {
	j, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
	}

	w.Write(j)
}
