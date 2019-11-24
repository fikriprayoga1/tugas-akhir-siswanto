package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

//Object
type responseObject struct {
	Response string
}

type updateDataObject struct {
	Name        string
	Temperature string
	Humidity    string
	OldName     string
}

type readDataObject struct {
	Name        string
	Temperature string
	Humidity    string
	LED         string
}

//Function Helper
func initDatabase(database *sql.DB) *sql.Tx {
	tx, err2 := database.Begin()
	if err2 != nil {
		log.Println(err2)
	}

	stmt, err3 := tx.Prepare("CREATE TABLE IF NOT EXISTS sbmList (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, temperature TEXT, humidity TEXT)")
	if err3 != nil {
		log.Println(err3)
	}
	stmt.Exec()
	defer stmt.Close()

	return tx

}

func updateResponseParser(request *http.Request) *updateDataObject {
	body, err0 := ioutil.ReadAll(request.Body)
	if err0 != nil {
		log.Println(err0)
	}
	var m updateDataObject
	err1 := json.Unmarshal(body, &m)
	if err1 != nil {
		log.Println(err1)
	}

	return &m
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/createData", createDataHandler)
	mux.HandleFunc("/readData", readDataHandler)
	mux.HandleFunc("/updateData", updateDataHandler)
	mux.HandleFunc("/updateData2", updateDataHandler2)
	mux.HandleFunc("/updateData3", updateDataHandler3)
	mux.HandleFunc("/deleteData", deleteDataHandler)
	
	http.ListenAndServe(":1810", mux)
}

func updateDataHandler3(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(0)

	mName := r.FormValue("name")
	mTemperature := r.FormValue("temperature")
	mHumidity := r.FormValue("humidity")
	mOldName := r.FormValue("oldname")

	database, err0 := sql.Open("sqlite3", "./sbm.db")
	if err0 != nil {
		log.Println(err0)
	}
	tx := initDatabase(database)
	defer database.Close()
	defer tx.Commit()

	stmt, err1 := tx.Prepare("UPDATE sbmList SET name=?, temperature=?, humidity=? WHERE name=?")
	if err1 != nil {
		log.Println(err1)
	}
	stmt.Exec(mName, mTemperature, mHumidity,mOldName)
	defer stmt.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	m2 := responseObject{"Create data success"}
	b, err2 := json.Marshal(m2)
	if err2 != nil {
		log.Println(err2)
	}
	w.Write(b)

}

func deleteDataHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(0)

	mName := r.FormValue("name")

	database, err0 := sql.Open("sqlite3", "./sbm.db")
	if err0 != nil {
		log.Println(err0)
	}
	tx := initDatabase(database)
	defer database.Close()
	defer tx.Commit()

	stmt, err0 := tx.Prepare("DELETE FROM sbmList WHERE name=?")
	if err0 != nil {
		log.Println(err0)
	}
	stmt.Exec(mName)
	defer stmt.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	m2 := responseObject{"delete data success"}
	b, err1 := json.Marshal(m2)
	if err1 != nil {
		log.Println(err1)
	}
	w.Write(b)
}

func createDataHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(0)

	mName := r.FormValue("name")
	mTemperature := r.FormValue("temperature")
	mHumidity := r.FormValue("humidity")
	log.Println(mName)
	log.Println(mTemperature)
	log.Println(mHumidity)

	database, err0 := sql.Open("sqlite3", "./sbm.db")
	if err0 != nil {
		log.Println(err0)
	}
	tx := initDatabase(database)
	defer database.Close()
	defer tx.Commit()

	stmt, err1 := tx.Prepare("INSERT INTO sbmList (name, temperature, humidity) VALUES (?, ?, ?)")
	if err1 != nil {
		log.Println(err1)
	}
	stmt.Exec(mName, mTemperature, mHumidity)
	defer stmt.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	m2 := responseObject{"Create data success"}
	b, err2 := json.Marshal(m2)
	if err2 != nil {
		log.Println(err2)
	}
	w.Write(b)

}

func readDataHandler(w http.ResponseWriter, r *http.Request) {
	database, err0 := sql.Open("sqlite3", "./sbm.db")
	if err0 != nil {
		log.Println(err0)
	}
	tx := initDatabase(database)
	defer database.Close()
	defer tx.Commit()

	mName := ""
	mTemperature := ""
	mHumidity := ""
	var mDeviceDataList []readDataObject
	mLED := "Mati"
	rows, err1 := tx.Query("SELECT name, temperature, humidity FROM sbmList")
	if err1 != nil {
		log.Println(err1)
	}
	for rows.Next() {
		rows.Scan(&mName, &mTemperature, &mHumidity)
		intTemp, err2 := strconv.ParseFloat(mTemperature, 64)
		if err2 != nil {
			log.Println(err2)
		}
		if intTemp >= 25 {
			mLED = "Hidup"

		}
		mDeviceDataList = append(mDeviceDataList, readDataObject{mName, mTemperature, mHumidity, mLED})

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	b, err2 := json.Marshal(mDeviceDataList)
	if err2 != nil {
		log.Println(err2)
	}
	w.Write(b)

}

func updateDataHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(0)

	mName := r.FormValue("name")
	mTemperature := r.FormValue("temperature")
	mHumidity := r.FormValue("humidity")
	mOldName := r.FormValue("oldName")

	database, err0 := sql.Open("sqlite3", "./sbm.db")
	if err0 != nil {
		log.Println(err0)
	}
	tx := initDatabase(database)
	defer database.Close()
	defer tx.Commit()

	stmt, err0 := tx.Prepare("UPDATE sbmList SET name=?, temperature=?, humidity=? WHERE name=?")
	if err0 != nil {
		log.Println(err0)
	}
	stmt.Exec(mName, mTemperature, mHumidity, mOldName)
	defer stmt.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	m2 := responseObject{"Update data success"}
	b, err1 := json.Marshal(m2)
	if err1 != nil {
		log.Println(err1)
	}
	w.Write(b)

}

func updateDataHandler2(w http.ResponseWriter, r *http.Request) {
	m := updateResponseParser(r)

	database, err0 := sql.Open("sqlite3", "./sbm.db")
	if err0 != nil {
		log.Println(err0)
	}
	tx := initDatabase(database)
	defer database.Close()
	defer tx.Commit()

	stmt, err0 := tx.Prepare("UPDATE sbmList SET name=?, temperature=?, humidity=? WHERE name=?")
	if err0 != nil {
		log.Println(err0)
	}
	stmt.Exec(m.Name, m.Temperature, m.Humidity, m.OldName)
	defer stmt.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	mLED := "M"
	intTemp, err1 := strconv.ParseFloat(m.Temperature, 64)
	if err1 != nil {
		log.Println(err1)
	}
	if intTemp >= 29 {
		mLED = "H"
	}
	w.Write([]byte(mLED))

}
