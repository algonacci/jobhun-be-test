package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Mahasiswa struct {
	Id                 int          `json:"id"`
	Nama               string       `json:"nama"`
	Usia               int          `json:"usia"`
	Gender             int          `json:"gender"`
	Tanggal_Registrasi sql.NullTime `json:"tanggal_registrasi"`
	Jurusan_Id         int          `json:"jurusan_id"`
	Hobi_Ids           []int        `json:"hobi_ids"`
}

type Jurusan struct {
	Id           int    `json:"id"`
	Nama_Jurusan string `json:"nama_jurusan"`
}

type Hobi struct {
	Id        int    `json:"id"`
	Nama_Hobi string `json:"nama_hobi"`
}

type MahasiswaResponse struct {
	Id                 int    `json:"id"`
	Nama               string `json:"nama"`
	Usia               int    `json:"usia"`
	Gender             int    `json:"gender"`
	Tanggal_Registrasi string `json:"tanggal_registrasi"`
	Nama_Jurusan       string `json:"nama_jurusan"`
	Hobi               string `json:"hobi"`
}

func main() {
	fmt.Println("🚀 [SERVER] is running on port http://localhost:8000")

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/jobhun_be_test")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", indexHandler)

	http.HandleFunc("/mahasiswa", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			createMahasiswa(db, w, r)
		case "GET":
			getAllMahasiswa(db, w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/mahasiswa/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getMahasiswaById(db, w, r)
		case "PUT":
			updateMahasiswa(db, w, r)
		case "DELETE":
			deleteMahasiswa(db, w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Fatal(http.ListenAndServe(":8000", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "API is running"}`)
}

func updateMahasiswa(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Parse ID Mahasiswa from URL parameter
	vars := r.URL.Query()
	id, err := strconv.Atoi(vars.Get("id"))
	if err != nil {
		http.Error(w, "Invalid Mahasiswa ID", http.StatusBadRequest)
		return
	}

	var m Mahasiswa
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Update Mahasiswa
	stmt, err := db.Prepare("UPDATE Mahasiswa SET Nama = ?, Usia = ?, Gender = ?, Tanggal_Registrasi = ? WHERE Id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := stmt.Exec(m.Nama, m.Usia, m.Gender, m.Tanggal_Registrasi, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Mahasiswa not found", http.StatusNotFound)
		return
	}

	// Update Jurusan
	stmt, err = db.Prepare("UPDATE Jurusan SET Nama_Jurusan = ? WHERE Id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err = stmt.Exec(m.Jurusan_Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rowsAffected, err = res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Jurusan not found", http.StatusNotFound)
		return
	}

	// Update Mahasiswa_Hobi
	_, err = db.Exec("DELETE FROM Mahasiswa_Hobi WHERE Id_Mahasiswa = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, h := range m.Hobi_Ids {
		stmt, err = db.Prepare("INSERT INTO Mahasiswa_Hobi (Id_Mahasiswa, Id_Hobi) VALUES (?, ?)")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res, err = stmt.Exec(id, h)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rowsAffected, err = res.RowsAffected()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if rowsAffected == 0 {
			http.Error(w, "Failed to insert Mahasiswa_Hobi", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Mahasiswa updated successfully")
}

func createMahasiswa(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var m Mahasiswa
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Insert Mahasiswa
	stmt, err := db.Prepare("INSERT INTO Mahasiswa (Id, Nama, Usia, Gender, Tanggal_Registrasi) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := stmt.Exec(m.Id, m.Nama, m.Usia, m.Gender, m.Tanggal_Registrasi)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Failed to insert Mahasiswa", http.StatusInternalServerError)
		return
	}

	// Insert Jurusan
	stmt, err = db.Prepare("INSERT INTO Jurusan (Id, Nama_Jurusan) VALUES (?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err = stmt.Exec(m.Jurusan_Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rowsAffected, err = res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Failed to insert Jurusan", http.StatusInternalServerError)
		return
	}

	// Insert Mahasiswa_Hobi
	for _, h := range m.Hobi_Ids {
		stmt, err = db.Prepare("INSERT INTO Mahasiswa_Hobi (Id_Mahasiswa, Id_Hobi) VALUES (?, ?)")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res, err = stmt.Exec(m.Id, h)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rowsAffected, err = res.RowsAffected()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if rowsAffected == 0 {
			http.Error(w, "Failed to insert Mahasiswa_Hobi", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Mahasiswa created successfully")
}

func getAllMahasiswa(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT m.Id, m.Nama, m.Usia, m.Gender, m.Tanggal_Registrasi, j.Nama_Jurusan, GROUP_CONCAT(h.Nama_Hobi SEPARATOR ', ') AS Hobi FROM Mahasiswa m LEFT JOIN Jurusan j ON m.Jurusan_Id = j.Id LEFT JOIN Mahasiswa_Hobi mh ON m.Id = mh.Id_Mahasiswa LEFT JOIN Hobi h ON mh.Id_Hobi = h.Id GROUP BY m.Id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var mahasiswas []MahasiswaResponse
	for rows.Next() {
		var m MahasiswaResponse
		if err := rows.Scan(&m.Id, &m.Nama, &m.Usia, &m.Gender, &m.Tanggal_Registrasi, &m.Nama_Jurusan, &m.Hobi); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		mahasiswas = append(mahasiswas, m)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mahasiswas)
}

func getMahasiswaById(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Parse ID Mahasiswa from URL parameter
	vars := r.URL.Query()
	id, err := strconv.Atoi(vars.Get("id"))
	if err != nil {
		http.Error(w, "Invalid Mahasiswa ID", http.StatusBadRequest)
		return
	}

	// Query Mahasiswa
	row := db.QueryRow("SELECT m.Id, m.Nama, m.Usia, m.Gender, m.Tanggal_Registrasi, j.Nama_Jurusan FROM Mahasiswa m LEFT JOIN Jurusan j ON m.Jurusan_Id = j.Id WHERE m.Id = ?", id)

	// Scan data into Mahasiswa struct
	var m Mahasiswa
	if err := row.Scan(&m.Id, &m.Nama, &m.Usia, &m.Gender, &m.Tanggal_Registrasi); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Mahasiswa not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Query Hobi
	rows, err := db.Query("SELECT h.Nama_Hobi FROM Mahasiswa_Hobi mh LEFT JOIN Hobi h ON mh.Id_Hobi = h.Id WHERE mh.Id_Mahasiswa = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var namaHobi string
		if err := rows.Scan(&namaHobi); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// m.Hobi = hobi

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

func deleteMahasiswa(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Parse ID Mahasiswa from URL parameter
	vars := r.URL.Query()
	id, err := strconv.Atoi(vars.Get("id"))
	if err != nil {
		http.Error(w, "Invalid Mahasiswa ID", http.StatusBadRequest)
		return
	}

	// Delete Mahasiswa from database
	res, err := db.Exec("DELETE FROM Mahasiswa WHERE Id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Mahasiswa not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Mahasiswa with ID %d has been deleted", id)
}
