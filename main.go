package main

import (
	"Project/connection"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"time"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()

	connection.DatabaseConnect()

	route.PathPrefix("/assets/").Handler(http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	route.HandleFunc("/", Home).Methods("GET")
	route.HandleFunc("/contact", Contact).Methods("GET")
	route.HandleFunc("/addMyProject", AddMyProject).Methods("GET")
	route.HandleFunc("/addMP", AddMP).Methods("POST")
	route.HandleFunc("/myProjectDetail/{id}", MyProjectDetail).Methods("GET")
	route.HandleFunc("/deleteMP/{id}", deleteMP).Methods("GET")
	route.HandleFunc("/editProject/{id}", edit).Methods("GET")
	route.HandleFunc("/update/{id}", update).Methods("POST")

	fmt.Println("Server Running")
	http.ListenAndServe("localhost:5000", route)

}

func Home(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/index.html")

	if err != nil {
		w.Write([]byte("error : " + err.Error()))
		return
	}

	data, _ := connection.Con.Query(context.Background(), "SELECT  id, title, description, duration FROM tb_project")
	fmt.Println(data)

	var result []MP
	for data.Next() {
		var each = MP{}

		var err = data.Scan(&each.ID, &each.Title, &each.Description, &each.Duration)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		result = append(result, each)
	}

	resData := map[string]interface{}{
		"MPs": result,
	}
	fmt.Println(result)

	tmpl.Execute(w, resData)

}
func Contact(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/contact.html")

	if err != nil {
		w.Write([]byte("error : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)

}
func AddMyProject(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/addMyProject.html")

	if err != nil {
		w.Write([]byte("error : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)

}

type MP struct {
	Title           string
	Description     string
	Duration        string
	ID              int
	StartDate       time.Time
	EndDate         time.Time
	Formatstartdate string
	Formatenddate   string
}

var dataMP = []MP{}

func AddMP(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var title = r.PostForm.Get("title")
	var description = r.PostForm.Get("description")
	var startDate = r.PostForm.Get("startDate")
	var endDate = r.PostForm.Get("endDate")

	layout := "2006-01-02"
	parsingstartdate, _ := time.Parse(layout, startDate)
	parsingenddate, _ := time.Parse(layout, endDate)

	hours := parsingenddate.Sub(parsingstartdate).Hours()
	days := hours / 24

	var duration string
	if days > 0 {
		duration = strconv.FormatFloat(days, 'f', 0, 64) + " days"
	}

	_, err = connection.Con.Exec(context.Background(), "INSERT INTO tb_project(title, start_date, end_date, description, duration) VAlUES ($1, $2, $3, $4, $5)", title, parsingstartdate, parsingenddate, description, duration)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}

func MyProjectDetail(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/myProjectDetail.html")

	if err != nil {
		w.Write([]byte("error : " + err.Error()))
		return
	}

	var MPDetail = MP{}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	err = connection.Con.QueryRow(context.Background(), "SELECT id, title, start_date, end_date, description, duration FROM tb_project WHERE id=$1", id).Scan(
		&MPDetail.ID, &MPDetail.Title, &MPDetail.StartDate, &MPDetail.EndDate, &MPDetail.Description, &MPDetail.Duration)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}

	MPDetail.Formatstartdate = MPDetail.StartDate.Format("2 january 2006")
	MPDetail.Formatenddate = MPDetail.EndDate.Format("2 january 2006")

	data := map[string]interface{}{
		"MP": MPDetail,
	}

	tmpl.Execute(w, data)

}

func deleteMP(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := connection.Con.Exec(context.Background(), "DELETE FROM tb_project WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func edit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset-utf8")

	var tmplt, err = template.ParseFiles("views/editProject.html")
	if err != nil {
		w.Write([]byte("file doesn't exist: " + err.Error()))
		return
	}
	var MPDetail = MP{}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	err = connection.Con.QueryRow(context.Background(), "SELECT id, title, start_date, end_date, description, duration FROM tb_project WHERE id=$1", id).Scan(&MPDetail.ID, &MPDetail.Title, &MPDetail.StartDate, &MPDetail.EndDate, &MPDetail.Description, &MPDetail.Duration)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}

	data := map[string]interface{}{
		"editProject": MPDetail,
	}
	tmplt.Execute(w, data)

}

func update(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var title = r.PostForm.Get("title")
	var description = r.PostForm.Get("description")
	var startDate = r.PostForm.Get("startDate")
	var endDate = r.PostForm.Get("endDate")

	layout := "2006-01-02"
	parsingstartdate, _ := time.Parse(layout, startDate)
	parsingenddate, _ := time.Parse(layout, endDate)

	hours := parsingenddate.Sub(parsingstartdate).Hours()
	days := hours / 24

	var duration string

	if days > 0 {
		duration = strconv.FormatFloat(days, 'f', 0, 64) + " days"
	}
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	sqlStatement := `UPDATE public.tb_project SET title=$2, start_date=$3, end_date=$4, description=$5, duration=$6
	WHERE id=$1;`

	_, err = connection.Con.Exec(context.Background(), sqlStatement, id, title, parsingstartdate, parsingenddate, description, duration)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
