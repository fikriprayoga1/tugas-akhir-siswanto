package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"

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
	Pressure    string
	Altitude    string
	LDR         string
	OldName     string
}

type readDataObject struct {
	Name        string
	Temperature string
	Humidity    string
	Pressure    string
	Altitude    string
	LDR         string
	LED         string
}

//Function Helper
var ledHolder = "Mati"

func initDatabase(database *sql.DB) *sql.Tx {
	tx, err2 := database.Begin()
	if err2 != nil {
		log.Println(err2)
	}

	stmt, err3 := tx.Prepare("CREATE TABLE IF NOT EXISTS sbmList (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, temperature TEXT, humidity TEXT, pressure TEXT,altitude TEXT, ldr TEXT)")
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
	http.HandleFunc("/createData", func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(0)

		mName := r.FormValue("name")
		mTemperature := r.FormValue("temperature")
		mHumidity := r.FormValue("humidity")
		mPressure := r.FormValue("pressure")
		mAltitude := r.FormValue("altitude")
		mLDR := r.FormValue("ldr")
		log.Println(mName)
		log.Println(mTemperature)
		log.Println(mHumidity)
		log.Println(mPressure)
		log.Println(mAltitude)
		log.Println(mLDR)

		database, err0 := sql.Open("sqlite3", "./sbm.db")
		if err0 != nil {
			log.Println(err0)
		}
		tx := initDatabase(database)
		defer database.Close()
		defer tx.Commit()

		stmt, err1 := tx.Prepare("INSERT INTO sbmList (name, temperature, humidity, pressure, altitude, ldr) VALUES (?, ?, ?, ?, ?, ?)")
		if err1 != nil {
			log.Println(err1)
		}
		stmt.Exec(mName, mTemperature, mHumidity, mPressure, mAltitude, mLDR)
		defer stmt.Close()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		m2 := responseObject{"Create data success"}
		b, err2 := json.Marshal(m2)
		if err2 != nil {
			log.Println(err2)
		}
		w.Write(b)

	})

	http.HandleFunc("/readData", func(w http.ResponseWriter, r *http.Request) {
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
		mPressure := ""
		mAltitude := ""
		mLDR := ""
		var mDeviceDataList []readDataObject
		rows, err1 := tx.Query("SELECT name, temperature, humidity, pressure, altitude, ldr  FROM sbmList")
		if err1 != nil {
			log.Println(err1)
		}
		for rows.Next() {
			rows.Scan(&mName, &mTemperature, &mHumidity, &mPressure, &mAltitude, &mLDR)

			mDeviceDataList = append(mDeviceDataList, readDataObject{mName, mTemperature, mHumidity, mPressure, mAltitude, mLDR, ledHolder})

		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		b, err2 := json.Marshal(mDeviceDataList)
		if err2 != nil {
			log.Println(err2)
		}
		w.Write(b)

	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var filepath = path.Join("views", "index.html")
		var tmpl, err0 = template.ParseFiles(filepath)
		if err0 != nil {
			http.Error(w, err0.Error(), http.StatusInternalServerError)
			return

		}

		// ------------------------------------------------------------------------------------------------------
		database, err1 := sql.Open("sqlite3", "./sbm.db")
		if err1 != nil {
			log.Println(err1)
		}
		tx := initDatabase(database)
		defer database.Close()
		defer tx.Commit()

		mName := ""
		mTemperature := ""
		mHumidity := ""
		mPressure := ""
		mAltitude := ""
		mLDR := ""
		rows, err2 := tx.Query("SELECT name, temperature, humidity, pressure, altitude, ldr  FROM sbmList")
		if err2 != nil {
			log.Println(err2)
		}
		for rows.Next() {
			rows.Scan(&mName, &mTemperature, &mHumidity, &mPressure, &mAltitude, &mLDR)

		}

		var data = map[string]interface{}{
			"temperature": mTemperature,
			"humidity":    mHumidity,
			"presure":     mPressure,
			"altitude":    mAltitude,
			"state":       ledHolder,
		}

		err2 = tmpl.Execute(w, data)
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		submit := r.FormValue("submit")
		log.Println(submit)
		if submit == "submit1" {
			ledHolder = "Hidup"
		} else {
			ledHolder = "Mati"
		}

		var filepath = path.Join("views", "index.html")
		var tmpl, err0 = template.ParseFiles(filepath)
		if err0 != nil {
			http.Error(w, err0.Error(), http.StatusInternalServerError)
			return

		}

		// ------------------------------------------------------------------------------------------------------
		database, err1 := sql.Open("sqlite3", "./sbm.db")
		if err1 != nil {
			log.Println(err1)
		}
		tx := initDatabase(database)
		defer database.Close()
		defer tx.Commit()

		mName := ""
		mTemperature := ""
		mHumidity := ""
		mPressure := ""
		mAltitude := ""
		mLDR := ""
		rows, err2 := tx.Query("SELECT name, temperature, humidity, pressure, altitude, ldr  FROM sbmList")
		if err2 != nil {
			log.Println(err2)
		}
		for rows.Next() {
			rows.Scan(&mName, &mTemperature, &mHumidity, &mPressure, &mAltitude, &mLDR)

		}

		var data = map[string]interface{}{
			"temperature": mTemperature,
			"humidity":    mHumidity,
			"presure":     mPressure,
			"altitude":    mAltitude,
			"state":       ledHolder,
		}

		err2 = tmpl.Execute(w, data)
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusInternalServerError)
		}

	})

	http.HandleFunc("/updateData3", func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(0)

		mName := r.FormValue("name")
		mTemperature := r.FormValue("temperature")
		mHumidity := r.FormValue("humidity")
		mPressure := r.FormValue("pressure")
		mAltitude := r.FormValue("altitude")
		mLDR := r.FormValue("ldr")
		mOldName := r.FormValue("oldname")

		database, err0 := sql.Open("sqlite3", "./sbm.db")
		if err0 != nil {
			log.Println(err0)
		}
		tx := initDatabase(database)
		defer database.Close()
		defer tx.Commit()

		stmt, err1 := tx.Prepare("UPDATE sbmList SET name=?, temperature=?, humidity=?, pressure=?, altitude=?, ldr=?  WHERE name=?")
		if err1 != nil {
			log.Println(err1)
		}
		stmt.Exec(mName, mTemperature, mHumidity, mPressure, mAltitude, mLDR, mOldName)
		defer stmt.Close()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		m2 := responseObject{"Create data success"}
		b, err2 := json.Marshal(m2)
		if err2 != nil {
			log.Println(err2)
		}
		w.Write(b)

	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("assets"))))

	log.Println("server started at localhost:1810")
	http.ListenAndServe(":1810", nil)
}
