package main

import (
	"fmt"
	"github.com/gorilla/mux"
	//"io/ioutil"
	"log"
	"net/http"
	//"net/url"
	//"strings"
	"github.com/mrsep18th/go_util/util_net"
	"html/template"
	//"os"
	"database/sql"
	"encoding/json"
	"github.com/mrsep18th/go_util/util_db"
	"sort"
)

var tmpl *template.Template

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/text_error", handler_text_error).Methods("POST")
	r.HandleFunc("/", handler_show_text_error).Methods("GET")
	r.HandleFunc("/get_records", handler_get_records).Methods("GET")
	r.HandleFunc("/filter", handler_filter).Methods("GET")

	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/")))
	r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")
	return r
}

func connect2Db(ipaddr string, port string) (sql.DB, error) {
	connString := util_db.GenerateDBConnectionString("root", "texterrcollector", "tcp", ipaddr, port, "text_err_collector")
	db, err := util_db.Connect2MySql(connString)
	return db, err
}

func main() {
	//args := os.Args[1:]
	db, err := connect2Db("127.0.0.1", "3306")

	if err != nil {
		log.Fatal(err)
	}

	tmpl = template.Must(template.ParseFiles("assets/index.html"))
	initStore(&db)

	r := newRouter()
	if err := http.ListenAndServe("192.168.2.106:8848", r); err != nil {
		log.Fatal(err)
	}
}

func handler_get_records(w http.ResponseWriter, r *http.Request) {
	records, err := store.GetRecords()

	if err != nil {
		log.Fatal(err)
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].Update_time > records[j].Update_time
	})

	jsonRecords, jsonErr := json.Marshal(records)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	w.Write((jsonRecords))
}

func handler_show_text_error(w http.ResponseWriter, r *http.Request) {
	err := tmpl.Execute(w, nil)

	if err != nil {
		fmt.Println(err)
	}
}

func handler_filter(w http.ResponseWriter, r *http.Request) {
	result_array, err := util_net.GetQueryValues(r, "text_content", "lang")

	if err != nil {
		log.Println(err)
	}

	text_content, lang := string(result_array[0]), string(result_array[1])

	log.Println("Filter : ", text_content, lang)
	store.AddFilter(text_content, lang)
	http.Redirect(w, r, "/", http.StatusFound)
}

func handler_text_error(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println("ParseForm() err: %v", err)
		return
	}

	data := r.PostForm["data"]

	newRecord := Record{}

	data_map := map[string]interface{}{}

	json.Unmarshal(([]byte)(data[0]), &data_map)

	newRecord.Text_content = data_map["text_content"].(string)
	newRecord.Language = data_map["language"].(string)
	newRecord.Update_time = int(data_map["update_time"].(float64))
	newRecord.User_name = "unknown"
	newRecord.File_name = ""
	newRecord.Version_str = "unknown"

	if data_map["user_name"] != nil {
		newRecord.User_name = data_map["user_name"].(string)
	}

	if data_map["file_name"] != nil {
		newRecord.File_name = data_map["file_name"].(string)
	}

	if data_map["version_str"] != nil {
		newRecord.Version_str = data_map["version_str"].(string)
	}

	err := store.UpdateRecord(&newRecord)
	if err != nil {
		log.Println(err)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
