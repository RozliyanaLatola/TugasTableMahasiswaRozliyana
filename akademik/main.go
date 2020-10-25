package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"fmt"
	"net/http"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"

)
var db *sql.DB
var err error
type yamlconfig struct {
	Connection struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Password string `yaml:"password"`
		User     string `yaml:"user"`
		Database string `yaml:"database"`
	}
}
type mahasiswa struct {
	Idmahasiswa string `json:"id_mahasiswa"`
	Nama        string `json:"nama"`
	Alamat      struct {
		Jalan     string `json:"jalan"`
		Kelurahan string `json:"kelurahan"`
		Kecamatan string `json:"kecamatan"`
		Kabupaten string `json:"kabupaten"`
		Provinsi  string `json:"provinsi"`
	} `json:"alamat"`
	Fakultas    string        `json:"fakultas"`
	Jurusan     string        `json:"jurusan"`
	Nilai []nilai `json:"Nilai"`
}

type nilai struct {
	Idmahasiswa  string  `json:"id_mahasiswa"`
	Idmatkul string  `json:"id_matkul"`
	Mkuliah   string  `json:"m_kuliah"`
	Nilai        float32 `json:"nilai"`
	Semester     int8    `json:"semester"`
}
func getNilai(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var mhsP []mahasiswa

	params := mux.Vars(r)

	sql := `SELECT
				id_mahasiswa,
				IFNULL(nama,'') nama,
				IFNULL(jalan,'') jalan,
				IFNULL(kelurahan,'') kelurahan,
				IFNULL(kecamatan,'') kecamatan,
				IFNULL(kabupaten,'') kabupaten,
				IFNULL(provinsi,'') provinsi,
				IFNULL(fakultas,'') fakultas,
				IFNULL(jurusan,'') jurusan				
			FROM mahasiswa WHERE id_mahasiswa IN (?)`

	result, err := db.Query(sql, params["id"])

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		var mhs mahasiswa
		err := result.Scan(&mhs.Idmahasiswa, &mhs.Nama, &mhs.Alamat.Jalan, &mhs.Alamat.Kelurahan, &mhs.Alamat.Kecamatan, &mhs.Alamat.Kabupaten, &mhs.Alamat.Provinsi, &mhs.Fakultas, &mhs.Jurusan)

		if err != nil {
			panic(err.Error())
		}

		sqlNilai := `SELECT
						id_mahasiswa		
						, mata_kuliah.id_matkul
						, mata_kuliah.m_kuliah
						, nilai
						, semester
					FROM
						nilai INNER JOIN mata_kuliah
							ON (nilai.id_matkul = mata_kuliah.id_matkul)
					WHERE nilai.id_mahasiswa = ?`

		resultNilai, errNilai := db.Query(sqlNilai, mhs.Idmahasiswa)

		defer resultNilai.Close()

		if errNilai != nil {
			panic(err.Error())
		}

		for resultNilai.Next() {
			var nilaiP nilai
			err := resultNilai.Scan(&nilaiP.Idmahasiswa, &nilaiP.Idmatkul, &nilaiP.Mkuliah, &nilaiP.Nilai, &nilaiP.Semester)
			if err != nil {
				panic(err.Error())
			}
			mhs.Nilai = append(mhs.Nilai, nilaiP)
		}
		mhsP = append(mhsP, mhs)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mhsP)

}
func getNilaiAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var mhsG []mahasiswa

	sql := `SELECT
				id_mahasiswa,
				IFNULL(nama,'') nama,
				IFNULL(jalan,'') jalan,
				IFNULL(kelurahan,'') kelurahan,
				IFNULL(kecamatan,'') kecamatan,
				IFNULL(kabupaten,'') kabupaten,
				IFNULL(provinsi,'') provinsi,
				IFNULL(fakultas,'') fakultas,
				IFNULL(jurusan,'') jurusan				
			FROM mahasiswa`

	result, err := db.Query(sql)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		var mhs2 mahasiswa
		err := result.Scan(&mhs2.Idmahasiswa, &mhs2.Nama, &mhs2.Alamat.Jalan, &mhs2.Alamat.Kelurahan, &mhs2.Alamat.Kecamatan, &mhs2.Alamat.Kabupaten, &mhs2.Alamat.Provinsi, &mhs2.Fakultas, &mhs2.Jurusan)

		if err != nil {
			panic(err.Error())
		}

		sqlNilai := `SELECT
						id_mahasiswa		
						, mata_kuliah.id_matkul
						, mata_kuliah.m_kuliah
						, nilai
						, semester
					FROM
						nilai INNER JOIN mata_kuliah
							ON (nilai.id_matkul = mata_kuliah.id_matkul)
					WHERE nilai.id_mahasiswa = ?`

		resultNilai, errNilai := db.Query(sqlNilai, mhs2.Idmahasiswa)

		defer resultNilai.Close()

		if errNilai != nil {
			panic(err.Error())
		}

		for resultNilai.Next() {
			var nilaiG nilai
			err := resultNilai.Scan(&nilaiG.Idmahasiswa, &nilaiG.Idmatkul, &nilaiG.Mkuliah, &nilaiG.Nilai, &nilaiG.Semester)
			if err != nil {
				panic(err.Error())
			}
			mhs2.Nilai = append(mhs2.Nilai, nilaiG)
		}
		mhsG = append(mhsG, mhs2)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mhsG)
}
func updateMahasiswa(w http.ResponseWriter, r *http.Request) {

	if r.Method == "PUT" {

		params := mux.Vars(r)

		newNama := r.FormValue("nama")
		newJalan := r.FormValue("jalan")
		newKelurahan := r.FormValue("kelurahan")
		newKecamatan := r.FormValue("kecamatan")
		newKabupaten := r.FormValue("kabupaten")
		newProvinsi := r.FormValue("provinsi")
		newFakultas := r.FormValue("fakultas")
		newJurusan := r.FormValue("jurusan")

		stmt, err := db.Prepare("UPDATE mahasiswa SET nama = ?, jalan = ?, kelurahan = ?, kecamatan = ?, kabupaten = ?, provinsi = ?, fakultas = ?, jurusan = ? WHERE id_mahasiswa = ?")

		_, err = stmt.Exec(newNama,newJalan,newKelurahan,newKecamatan,newKabupaten,newProvinsi,newFakultas,newJurusan, params["id"])

		if err != nil {
			fmt.Fprintf(w, "Data not found or Request error")
		}

		fmt.Fprintf(w, "Mahasiswa with id_mahasiswa = %s was updated", params["id"])
	}
}
func createMahasiswa(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		Idmahasiswa := r.FormValue("id_mahasiswa")
		Nama := r.FormValue("nama")
		Jalan := r.FormValue("jalan")
		Kelurahan := r.FormValue("kelurahan")
		Kecamatan := r.FormValue("kecamatan")
		Kabupaten := r.FormValue("kabupaten")
		Provinsi := r.FormValue("provinsi")
		Fakultas := r.FormValue("fakultas")
		Jurusan := r.FormValue("jurusan")

		stmt, err := db.Prepare("INSERT INTO mahasiswa (id_mahasiswa, nama, jalan, kelurahan, kecamatan, kabupaten, provinsi, fakultas, jurusan) VALUES (?,?,?,?,?,?,?,?,?)")

		_, err = stmt.Exec(Idmahasiswa, Nama, Jalan, Kelurahan, Kecamatan, Kabupaten, Provinsi, Fakultas, Jurusan)

		if err != nil {
			fmt.Fprintf(w, "Data Duplicate")
		} else {
			fmt.Fprintf(w, "Data Created")
		}

	}
}
func main() {
	yamlFile, err := ioutil.ReadFile("../Yaml/config.yml")
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return
	}
	var yamlConfig yamlconfig
	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
	}

	host := yamlConfig.Connection.Host
	port := yamlConfig.Connection.Port
	user := yamlConfig.Connection.User
	pass := yamlConfig.Connection.Password
	data := yamlConfig.Connection.Database

	var (
		mySQL = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", user, pass, host, port, data)
	)

	db, err = sql.Open("mysql", mySQL)
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	// Init router
	r := mux.NewRouter()

	// Route handles & endpoints
	r.HandleFunc("/mahasiswaG", getNilaiAll).Methods("GET")
	r.HandleFunc("/mahasiswa/{id}", updateMahasiswa).Methods("PUT")
	r.HandleFunc("/mahasiswaT", createMahasiswa).Methods("POST")
	r.HandleFunc("/mahasiswaP/{id}", getNilai).Methods("GET")

	fmt.Println("Server on :8181")
	// Start server
	log.Fatal(http.ListenAndServe(":8181", r))
}