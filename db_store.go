package main

import (
	"database/sql"
	"log"
	"sync"
)

type db_store struct {
	db  *sql.DB
	mux sync.Mutex
}

type Record struct {
	Text_content string
	Language     string
	Update_time  int
	User_name    string
	File_name    string
	Version_str  string
	Ignore       int
}

type Store interface {
	UpdateRecord(rec *Record) error
	GetRecords() ([]*Record, error)
	AddFilter(txtContent string, lang string) error
}

var store Store

func initStore(dbInst *sql.DB) {
	store = &db_store{db: dbInst}
}

func (dbs *db_store) AddFilter(txtContent string, lang string) error {
	dbs.mux.Lock()
	defer dbs.mux.Unlock()
	rows, err := dbs.db.Query("UPDATE errors SET ignore_flag = 1 Where text_content=? AND language=?", txtContent, lang)

	if err != nil {
		log.Println("AddFilter: ", err)
		return err
	}

	defer rows.Close()
	return err
}

func (dbs *db_store) UpdateRecord(rec *Record) error {
	dbs.mux.Lock()
	defer dbs.mux.Unlock()
	if dbs.CheckExistence(rec.Text_content, rec.Language) {
		rows, err := dbs.db.Query("UPDATE errors SET update_time=?, user_name = ?, file_name = ?, version_str = ? Where text_content=? AND language=?", rec.Update_time, rec.User_name, rec.File_name, rec.Version_str, rec.Text_content, rec.Language)
		if err == nil {
			rows.Close()
		}

		return err
	} else {
		rows, err := dbs.db.Query("INSERT INTO errors (text_content, language, update_time, user_name, file_name, version_str) VALUES(?,?,?,?,?,?)", rec.Text_content, rec.Language, rec.Update_time, rec.User_name, rec.File_name, rec.Version_str)
		if err == nil {
			rows.Close()
		}

		return err
	}

}

func (dbs *db_store) CheckExistence(text_content string, language string) bool {
	rows, err := dbs.db.Query("select 1 from errors where text_content=? AND language=? limit 1;", text_content, language)

	if err != nil {
		log.Fatal("CheckExistence: ", err)
	}

	defer rows.Close()

	rows.Next()

	var result int
	rows.Scan(&result)

	//fmt.Println("CheckExistence = ", result)

	return result == 1
}

func (dbs *db_store) GetRecords() ([]*Record, error) {
	rows, err := dbs.db.Query("SELECT text_content, language, update_time, user_name, file_name, version_str, ignore_flag from errors")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := []*Record{}

	for rows.Next() {
		record := &Record{}
		rows.Scan(&record.Text_content, &record.Language, &record.Update_time, &record.User_name, &record.File_name, &record.Version_str, &record.Ignore)
		records = append(records, record)
	}
	return records, nil
}
