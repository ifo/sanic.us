package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ifo/sanic"
)

var port = flag.Int("port", 3000, "Port to run the server on")

func main() {
	flag.Parse()

	http.Handle("/", router())

	log.Printf("Starting server on port %d\n", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/{id:[7-9]|10}", idHandler).Methods("GET")
	// TODO handle .html, .text, and .json appends
	return r
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/10", http.StatusTemporaryRedirect)
}

var indexPage = `<!DOCTYPE HTML>
<html>
  <head>
    <title>id.sanic.us</title>
    <style>
      html, body {
        height: 100%;
        margin: 0;
        display: flex;
        align-items: center;
        justify-content: center;
      }

      #id {
        font-size: 3em;
        border: none;
        text-align: center;
      }
    </style>
  </head>
  <body>
    <div class='container'>
      <input id='id' type='text' value='{{.}}' readonly autofocus />
    </div>
    <script>
      document.addEventListener('DOMContentLoaded', function() {
        var id = document.getElementById('id');
        id.select();
      });
    </script>
  </body>
</html>
`

func idHandler(w http.ResponseWriter, r *http.Request) {
	length, _ := strconv.Atoi(mux.Vars(r)["id"])
	worker, _ := generateXLengthWorker(length)
	id := worker.IDString(worker.NextID())

	tmpl := template.Must(template.New("index").Parse(indexPage))
	tmpl.Execute(w, id)
}

func generateXLengthWorker(x int) (*sanic.Worker, error) {
	switch x {
	case 7:
		return &sanic.SevenLengthWorker, nil
	case 8:
		return &sanic.EightLengthWorker, nil
	case 9:
		return &sanic.NineLengthWorker, nil
	case 10:
		return &sanic.TenLengthWorker, nil
	default:
		return nil, fmt.Errorf("%d is not a number between 7 and 10")
	}
}
