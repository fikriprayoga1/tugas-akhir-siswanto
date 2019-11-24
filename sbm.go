package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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

func updateDataAndroid(aName string, aTemperature string, aHumidity string, aOldName string) {
	isDone := true
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
	rows, err1 := tx.Query("SELECT name, temperature, humidity FROM sbmList")
	if err1 != nil {
		log.Println(err1)
	}
	for rows.Next() {
		rows.Scan(&mName, &mTemperature, &mHumidity)
		if mName == aName {
			if ((aName) != (aOldName)) || ((mTemperature) != (aTemperature)) || ((mHumidity) != (aHumidity)) {
				if isDone {
					client := &http.Client{}
					postData := []byte("{\"to\": \"/topics/update\", \"data\": {\"message\": \"Server data is updated\"}}")
					req, err := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", bytes.NewReader(postData))
					if err != nil {
						os.Exit(1)
					}
					req.Header.Add("Content-Type", "application/json")
					req.Header.Add("Authorization", "key=AIzaSyBjpyC3bPaCekMJy81Irf1TxsZAe7CYRP4")
					resp, err := client.Do(req)
					defer resp.Body.Close()

					isDone = false

				}

			}

		}

	}

}

func updateDataAndroid2() {
	client := &http.Client{}
	postData := []byte("{\"to\": \"/topics/update\", \"data\": {\"message\": \"Server data is updated\"}}")
	req, err := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", bytes.NewReader(postData))
	if err != nil {
		os.Exit(1)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "key=AIzaSyBjpyC3bPaCekMJy81Irf1TxsZAe7CYRP4")
	resp, err := client.Do(req)
	defer resp.Body.Close()

}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/createData", createDataHandler)
	mux.HandleFunc("/readData", readDataHandler)
	mux.HandleFunc("/updateData", updateDataHandler)
	mux.HandleFunc("/updateData2", updateDataHandler2)
	mux.HandleFunc("/deleteData", deleteDataHandler)

	http.ListenAndServe(":1810", mux)
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

	updateDataAndroid2()

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
	rows, err1 := tx.Query("SELECT name, temperature, humidity FROM sbmList")
	if err1 != nil {
		log.Println(err1)
	}
	for rows.Next() {
		rows.Scan(&mName, &mTemperature, &mHumidity)
		mDeviceDataList = append(mDeviceDataList, readDataObject{mName, mTemperature, mHumidity})

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

	updateDataAndroid(mName, mTemperature, mHumidity, mOldName)
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

	m2 := responseObject{"Update data success"}
	b, err1 := json.Marshal(m2)
	if err1 != nil {
		log.Println(err1)
	}
	w.Write(b)

	updateDataAndroid(m.Name, m.Temperature, m.Humidity, m.OldName)

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

	m2 := responseObject{"Delete data success"}
	b, err1 := json.Marshal(m2)
	if err1 != nil {
		log.Println(err1)
	}
	w.Write(b)

	updateDataAndroid2()

}
