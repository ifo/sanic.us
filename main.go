package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ifo/sanic"
)

var port = flag.Int("port", 3000, "Port to run the server on")

type Context struct {
	WorkerMap map[string]*sanic.Worker
	Vars      map[string]string
	Templates *template.Template
}

func main() {
	flag.Parse()

	c := Context{
		WorkerMap: generateWorkers(),
		Templates: template.Must(template.New("id").Parse(indexPage)),
	}

	http.Handle("/", router(c))

	log.Printf("Starting server on port %d\n", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request, c Context) {
	http.Redirect(w, r, "/10", http.StatusTemporaryRedirect)
}

func idHandler(w http.ResponseWriter, r *http.Request, c Context) {
	wid := c.Vars["id"]
	worker := c.WorkerMap[wid]
	id := worker.IDString(worker.NextID())
	tmpl := c.Templates.Lookup("id")
	tmpl.Execute(w, id)
}

func injectContext(fn func(w http.ResponseWriter, r *http.Request, c Context), c Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.Vars = mux.Vars(r)
		fn(w, r, c)
	}
}

func router(c Context) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", injectContext(indexHandler, c)).Methods("GET")
	r.HandleFunc("/{id:[7-9]|10}", injectContext(idHandler, c)).Methods("GET")
	// TODO handle .html, .text, and .json appends
	return r
}

func generateWorkers() map[string]*sanic.Worker {
	return map[string]*sanic.Worker{
		"10": sanic.NewWorker10(0),
		"9":  sanic.NewWorker9(0),
		"8":  sanic.NewWorker8(),
		"7":  sanic.NewWorker7(),
	}
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
