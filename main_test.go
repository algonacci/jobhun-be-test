package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateMahasiswaEndpoint(t *testing.T) {
	// Setup test database
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/jobhun_be_test")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			createMahasiswa(db, w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	defer ts.Close()

	// Prepare request body
	requestBody := []byte(`{"nama": "John Doe", "usia": 25, "gender": 1, "tanggal_registrasi": "2023-04-20T11:00:00Z", "jurusan_id": 1, "hobi_id": [1, 2]}`)

	// Send HTTP POST request
	res, err := http.Post(ts.URL+"/mahasiswa", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Check response status code
	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %d but got %d", http.StatusCreated, res.StatusCode)
	}

	// Check response body
	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	var response map[string]interface{}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		t.Fatal(err)
	}
	if _, ok := response["id"]; !ok {
		t.Error("Expected response body to contain 'id' field")
	}
}

func TestGetAllMahasiswaEndpoint(t *testing.T) {
	// Setup test database
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/jobhun_be_test")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getAllMahasiswa(db, w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	defer ts.Close()

	// Send HTTP GET request
	res, err := http.Get(ts.URL + "/mahasiswa")
	if err != nil {
		t.Fatal(err)
	}

	// Check response status code
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, res.StatusCode)
	}

	// Check response body
	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	var response []map[string]interface{}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		t.Fatal(err)
	}
	if len(response) == 0 {
		t.Error("Expected response body to contain data")
	}
}
