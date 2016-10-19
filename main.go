package main

import (
	"fmt"
	"github.com/andlabs/ui"
	_ "github.com/jmoiron/sqlx"
	"github.com/skratchdot/open-golang/open"
	"net"
	"net/http"
	"time"
)

func main() {

	http.Handle("/", http.FileServer(http.Dir("static")))

	listener, err := net.Listen("tcp", ":8080")
	go http.Serve(listener, nil)

	err = open.Run("http://localhost:8080/")
	if err != nil {
		fmt.Println(err)
	}

	err = ui.Main(func() {
		window := ui.NewWindow("Dizionario", 300, 300, false)
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
	})
	if err != nil {
		panic(err)
	}
}
