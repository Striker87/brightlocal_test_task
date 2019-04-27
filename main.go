package main

import (
	"./check"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

type Data struct {
	Method string `json:"method"`
	Key    string `json:"key"`
	Value  string `json:"value,omitempty"`
	Error  string `json:"error,omitempty"`
	Result string `json:"result,omitempty"`
	sync.RWMutex
	storage map[string]string
}

var data Data

func main() {
	data = Data{storage: make(map[string]string)}

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
		w.Write([]byte(`{"error":"wrong json"}`))
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
	result := Data{}

	if check.KeyExists(d.Key, d.storage) {
		result = Data{
			Method: "GET",
			Key:    d.Key,
			Value:  d.storage[d.Key],
		}
	} else {
		result = Data{
			Method: "GET",
			Key:    d.Key,
			Error:  "not found",
		}
	}

	render(result, w)
}

func (d *Data) set(w http.ResponseWriter) {
	result := Data{}

	if check.Value(d.Value) {
		result = Data{
			Method: "SET",
			Value:  d.Value,
			Error:  "value too long",
		}
	} else if check.Key(d.Key) {
		result = Data{
			Method: "SET",
			Key:    d.Key,
			Error:  "key too long or empty",
		}
	} else if check.StorageSize(d.storage) {
		result = Data{
			Method: "SET",
			Key:    d.Key,
			Error:  "storage is full. maximum size 1024 keys",
		}
	} else {
		d.Lock()
		d.storage[d.Key] = d.Value
		d.Unlock()

		result = Data{
			Method: "SET",
			Key:    d.Key,
			Value:  d.Value,
		}
	}

	render(result, w)
}

func (d *Data) exists(w http.ResponseWriter) {
	result := Data{}

	if check.KeyExists(d.Key, d.storage) {
		result = Data{
			Method: "EXISTS",
			Key:    d.Key,
			Result: "exists",
		}
	} else {
		result = Data{
			Method: "EXISTS",
			Key:    d.Key,
			Result: "not exists",
		}
	}

	render(result, w)
}

func (d *Data) remove(w http.ResponseWriter) {
	result := Data{}

	if check.KeyExists(d.Key, d.storage) {
		delete(d.storage, d.Key)

		result = Data{
			Method: "DELETE",
			Key:    d.Key,
			Result: "success",
		}
	} else {
		result = Data{
			Method: "DELETE",
			Key:    d.Key,
			Error:  "not found",
		}
	}

	render(result, w)
}

func wrongMethod(w http.ResponseWriter) {
	methods := []string{"GET", "SET", "DELETE", "EXISTS"}

	result := Data{
		Error: "Method not allowed. Allowed methods is: " + strings.Join(methods, ","),
	}

	render(result, w)
}

func render(result Data, w http.ResponseWriter) {
	j, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
	}

	w.Write(j)
}
