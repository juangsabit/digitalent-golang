package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type task struct {
	Id       int
	Task     string
	Assignee string
	Deadline string
	Done     int
}

type response struct {
	Status bool
	Pesan  string
	Data   []task
}

func koneksi() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/go")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func tampil(pesan string) response {
	db, err := koneksi()
	if err != nil {
		return response{
			Status: false,
			Pesan:  "Gagal koneksi " + err.Error(),
			Data:   []task{},
		}
	}
	defer db.Close()
	dataTask, err := db.Query("select * from task")
	if err != nil {
		return response{
			Status: false,
			Pesan:  "Gagal query " + err.Error(),
			Data:   []task{},
		}
	}
	defer dataTask.Close()
	var hasil []task
	for dataTask.Next() {
		var tsk = task{}
		var err = dataTask.Scan(&tsk.Id, &tsk.Task, &tsk.Assignee, &tsk.Deadline, &tsk.Done)
		if err != nil {
			return response{
				Status: false,
				Pesan:  "Gagal baca " + err.Error(),
				Data:   []task{},
			}
		}
		hasil = append(hasil, tsk)
	}
	if err != nil {
		return response{
			Status: false,
			Pesan:  "Error " + err.Error(),
			Data:   []task{},
		}
	}
	return response{
		Status: true,
		Pesan:  pesan,
		Data:   hasil,
	}
}

func gettsk(Id int) response {
	db, err := koneksi()
	if err != nil {
		return response{
			Status: false,
			Pesan:  "Gagal koneksi " + err.Error(),
			Data:   []task{},
		}
	}
	defer db.Close()
	dataTask, err := db.Query("select * from task where Id=?", Id)
	if err != nil {
		return response{
			Status: false,
			Pesan:  "Gagal query " + err.Error(),
			Data:   []task{},
		}
	}
	defer dataTask.Close()
	var hasil []task
	for dataTask.Next() {
		var tsk = task{}
		var err = dataTask.Scan(&tsk.Id, &tsk.Task, &tsk.Assignee, &tsk.Deadline, &tsk.Done)
		if err != nil {
			return response{
				Status: false,
				Pesan:  "Gagal baca " + err.Error(),
				Data:   []task{},
			}
		}
		hasil = append(hasil, tsk)
	}
	if err != nil {
		return response{
			Status: false,
			Pesan:  "Error " + err.Error(),
			Data:   []task{},
		}
	}
	return response{
		Status: true,
		Pesan:  "Berhasil tampil",
		Data:   hasil,
	}
}

func tambah(Id, Task string, Assignee string, Deadline string) response {
	db, err := koneksi()
	if err != nil {
		return response{
			Status: false,
			Pesan:  "Gagal koneksi " + err.Error(),
			Data:   []task{},
		}
	}
	defer db.Close()
	Done := 0
	_, err = db.Exec("insert into task values (?,?,?,?,?)", Id, Task, Assignee, Deadline, Done)
	if err != nil {
		return response{
			Status: false,
			Pesan:  "Gagal query insert " + err.Error(),
			Data:   []task{},
		}
	}
	return response{
		Status: true,
		Pesan:  "Berhasil tampil",
		Data:   []task{},
	}
}

func ubah(Id string, Task string, Assignee string, Deadline string) response {
	db, err := koneksi()
	if err != nil {
		return response{
			Status: false,
			Pesan:  "Gagal koneksi " + err.Error(),
			Data:   []task{},
		}
	}
	defer db.Close()
	_, err = db.Exec("update task set Task=?, Assignee=?, Deadline=? where Id=?", Task, Assignee, Deadline, Id)
	if err != nil {
		return response{
			Status: false,
			Pesan:  "Gagal query update " + err.Error(),
			Data:   []task{},
		}
	}
	return response{
		Status: true,
		Pesan:  "Berhasil tampil",
		Data:   []task{},
	}
}

func markdone(Id string) response {
	db, err := koneksi()
	if err != nil {
		return response{
			Status: false,
			Pesan:  "Gagal koneksi " + err.Error(),
			Data:   []task{},
		}
	}
	defer db.Close()
	_, err = db.Exec("update task set Done=1 where Id=?", Id)
	if err != nil {
		return response{
			Status: false,
			Pesan:  "Gagal query markdone " + err.Error(),
			Data:   []task{},
		}
	}
	return response{
		Status: true,
		Pesan:  "Berhasil tampil",
		Data:   []task{},
	}
}

func hapus(Id string) response {
	db, err := koneksi()
	if err != nil {
		return response{
			Status: false,
			Pesan:  "Gagal koneksi " + err.Error(),
			Data:   []task{},
		}
	}
	defer db.Close()
	_, err = db.Exec("delete task where Id=?", Id)
	if err != nil {
		return response{
			Status: false,
			Pesan:  "Gagal query delete " + err.Error(),
			Data:   []task{},
		}
	}
	return response{
		Status: true,
		Pesan:  "Berhasil tampil",
		Data:   []task{},
	}
}

func kontroler(w http.ResponseWriter, r *http.Request) {
	var tampilHtml, errTampil = template.ParseFiles("views/tampil.html")
	if errTampil != nil {
		fmt.Println(errTampil.Error())
		return
	}
	var tambahHtml, errTambah = template.ParseFiles("views/tambah.html")
	if errTambah != nil {
		fmt.Println(errTambah.Error())
		return
	}
	var ubahHtml, errUbah = template.ParseFiles("views/ubah.html")
	if errUbah != nil {
		fmt.Println(errUbah.Error())
		return
	}
	var hapusHtml, errHapus = template.ParseFiles("views/hapus.html")
	if errHapus != nil {
		fmt.Println(errHapus.Error())
		return
	}

	switch r.Method {
	case "GET":
		aksi := r.URL.Query()["aksi"]
		if len(aksi) == 0 {
			tampilHtml.Execute(w, tampil("Berhasil tampil"))
		} else if aksi[0] == "tambah" {
			tambahHtml.Execute(w, nil)
		} else if aksi[0] == "ubah" {
			strId := r.URL.Query()["id"]
			intVar, _ := strconv.Atoi(strId[0])
			ubahHtml.Execute(w, gettsk(intVar))
		} else if aksi[0] == "markdone" {
			strId := r.URL.Query()["id"]
			hasil := markdone(strId[0])
			tampilHtml.Execute(w, tampil(hasil.Pesan))
		} else if aksi[0] == "hapus" {
			strId := r.URL.Query()["id"]
			intVar, _ := strconv.Atoi(strId[0])
			hapusHtml.Execute(w, gettsk(intVar))
		} else {
			tampilHtml.Execute(w, tampil("Berhasil tampil"))
		}
	case "POST":
		err := r.ParseForm()
		if err != nil {
			fmt.Println(w, "Error :", err)
			return
		}
		Id := r.FormValue("id")
		Task := r.FormValue("task")
		Assignee := r.FormValue("assignee")
		Deadline := r.FormValue("deadline")
		aksi := r.URL.Path
		// println(Task, Assignee, Deadline, aksi)
		if aksi == "/tambah" {
			hasil := tambah(Id, Task, Assignee, Deadline)
			tampilHtml.Execute(w, tampil(hasil.Pesan))
			aksi = "tampil"
		} else if aksi == "/ubah" {
			hasil := ubah(Id, Task, Assignee, Deadline)
			tampilHtml.Execute(w, tampil(hasil.Pesan))
		} else if aksi == "/hapus" {
			hasil := hapus(Id)
			tampilHtml.Execute(w, tampil(hasil.Pesan))
		} else {
			tampilHtml.Execute(w, tampil("Berhasil tampil"))
		}
	default:
		fmt.Println(w, "Method error")
	}
}

func connect() (*sql.DB, error) {
	db, err := sql.Open("MySql", "root:@tcp(127:0.0.1:3306)/go")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/", kontroler)
	var address = "localhost:9000"
	fmt.Printf("server started at %s\n", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
