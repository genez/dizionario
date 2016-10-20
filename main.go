package main

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/skratchdot/open-golang/open"
	"net"
	"net/http"
	"log"
	"os"
)

type Termine struct {
	Termine     string `db:"termine" json:"termine"`
	Definizione string `db:"definizione" json:"definizione"`
}

var db *sqlx.DB

func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	fmt.Println(q)

	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Queryx("SELECT * FROM definizioni WHERE definizioni MATCH ? ORDER BY rank;", q)
	if (err != nil) {
		log.Fatal(err)
	}

	termini := make([]Termine, 0)
	if err == nil {
		for rows.Next() {
			var t Termine
			err = rows.StructScan(&t)
			termini = append(termini, t)
		}

	} else {
		log.Panic(err.Error())
	}
	json.NewEncoder(w).Encode(termini)
}

func main() {
	//crea il file db3 con questo comando (io uso SQLiteSpy):
	//CREATE VIRTUAL TABLE definizioni USING fts5(termine, definizione);
	//https://sqlite.org/fts5.html

	os.Remove("dizionario.db3")

	db = sqlx.MustConnect("sqlite3", "file:dizionario.db3?loc=auto")
	defer db.Close()

	db.MustExec("CREATE VIRTUAL TABLE definizioni USING fts5(termine, definizione);")
	db.MustExec("INSERT INTO definizioni VALUES ('cavallo','animale equino molto bello')")
	db.MustExec("INSERT INTO definizioni VALUES ('sviluppatore','animale umano un po meno bello del cavallo')")
	db.MustExec("INSERT INTO definizioni VALUES ('cliente','essere mitologico metà umano e metà cavallo')")
	db.MustExec("INSERT INTO definizioni VALUES ('delfino','animale curioso e sapiente. prima di andarsene ringrazia per tutto il pesce')")

	http.HandleFunc("/search", searchHandler)
	http.Handle("/", http.FileServer(http.Dir("static")))

	listener, err := net.Listen("tcp", ":8080")
	go http.Serve(listener, nil)

	err = open.Run("http://localhost:8080/")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Chiudi questa finestra (o premi CTRL+c) per uscire...")

	select {}
}
