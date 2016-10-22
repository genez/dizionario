# dizionario

commenti sul file main.go

la parte `import` serve per referenziare le librerie esterne (in Go si chiamano package)

`var db *sqlx.DB` -> definisco un puntatore "globale" al database, può essere visto e acceduto in tutto il programma
 
passiamo alla func `main()` che rappresenta il punto di ingresso del programma quando viene lanciato
 
`os.Remove("dizionario.db3")` -> cancella (se esiste già) il file dizionario.db3 che conterrà il database. faccio così per poterlo poi ricreare di nuovo, ma nella realtà il file sarà pre-popolato e oltretutto essendo su sopporto readonly (cd-rom) questa chiamata andrà in errore

`db = sqlx.MustConnect("sqlite3", "file:dizionario.db3?loc=auto")` -> chiedo al driver sqlite3 di connettersi al file dizionario.db3

`defer db.Close()` -> chiedo al runtime di Go di chiudere il db al termine della funzione main()

`db.MustExec...` -> creazione della tabella con il modulo FTS5 di full-text search e vari statement di insert per popolarla con dei dati a caso

`http.HandleFunc("/search", searchHandler)` -> chiedo alla libreria http di Go di chiamare la funzione 'searchHandler' quando l'url è del tipo http://server/search

`http.Handle("/", http.FileServer(http.Dir("static")))` -> chiedo alla libreria http di Go di servire file statici che si trovano nella cartella "static"


`listener, err := net.Listen("tcp", ":8080")` -> mi metto in ascolto sulla porta 8080 TCP
`go http.Serve(listener, nil)` -> chiedo alla libreria http di Go di servire richieste HTTP sul listener appena creato (8080 TCP)

`err = open.Run("http://localhost:8080/")` -> usando il package "open" (github.com/skratchdot/open-golang/open) lancio il browser di default sulla root del mio sito

`fmt.Println("Chiudi questa finestra (o premi CTRL+c) per uscire...")` -> scrivo qualcosa a video

`select {}` -> attendo indefinitamente



funzione searchHandler()
ha 2 parametri: r è la Request e w è la Response

`q := r.URL.Query().Get("q")` -> prendo il parametro q della querystring (ad esempio: http://localhost:8080/search?q=pippo mi ritorna pippo)

`w.Header().Set("Content-Type", "application/json")` -> indico al client che la risposta è di tipo JSON

`rows, err := db.Queryx("SELECT * FROM definizioni WHERE definizioni MATCH ? ORDER BY rank;", q)` -> estraggo dal db sqlite le definizioni che hanno un match con la parola chiave

`defer rows.Close()` -> chiudo il "cursore" sulle righe ritornate al termine della funzione

`termini := make([]Termine, 0)` -> creo una lista di termini con dimensione iniziale 0 (nessun termine)

poi per ogni riga
`var t Termine`
`err = rows.StructScan(&t)`
`termini = append(termini, t)`
creo un termine, lo valorizzo con i dati dal db e lo aggiungo alla lista

`json.NewEncoder(w).Encode(termini)` -> scrivo al client l'elenco di termini in formato JSON

that's it!
