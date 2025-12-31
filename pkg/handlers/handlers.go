package handlers

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"pastecode/pkg/app"
	"pastecode/pkg/paste"
)

const MaxBodySize = 2 << 20

func Index(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.Pastecodes.GC()

		tpl, err := template.ParseFiles("templates/base.html",
			"templates/head.html",
			"templates/topmenu.html",
			"templates/footer.html",
			"templates/index.html")
		if err != nil {
			log.Printf("error parsing template files: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = tpl.Execute(w, app.Pastecodes)
		if err != nil {
			log.Printf("error executing template files: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func Add(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			tpl, err := template.ParseFiles("templates/base.html",
				"templates/head.html",
				"templates/topmenu.html",
				"templates/footer.html",
				"templates/add.html")
			if err != nil {
				log.Printf("error parsing template files: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = tpl.Execute(w, "base.html")
			if err != nil {
				log.Printf("error executing template files: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		if r.Method == "POST" {
			r.Body = http.MaxBytesReader(w, r.Body, MaxBodySize)

			err := r.ParseForm()
			if err != nil {
				log.Printf("error form parsing: %s", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			p, err := paste.NewPastecode(r.FormValue("username"), r.FormValue("code"))
			if err != nil {
				log.Printf("error creating paste: %s", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			err = app.Pastecodes.Add(p)
			if err != nil {
				log.Printf("error adding paste: %s", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}

func Paste(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuidRequest := r.PathValue("uuid")
		id, err := paste.ParseUUID(uuidRequest)
		if err != nil {
			log.Printf("error parsing uuid from request: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		p, err := app.Pastecodes.FindPaste(id)
		if err != nil {
			log.Printf("error lookup uuid: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tpl, err := template.ParseFiles("templates/base.html",
			"templates/head.html",
			"templates/topmenu.html",
			"templates/footer.html",
			"templates/paste.html")
		if err != nil {
			log.Printf("error parsing template files: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = tpl.Execute(w, p)
		if err != nil {
			log.Printf("error executing template files: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func Static(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fs := http.FileServer(http.Dir("static"))
		http.StripPrefix("/static/", fs).ServeHTTP(w, r)
	}
}

func Healthz(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func Readyz(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

// LoggingMiddleware wraps every request with logging data.
func LoggingMiddleware(app *app.Application, next http.Handler) http.Handler {
	start := time.Now()
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)

			app.Sugar.Infow("Request processed",
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(start),
			)
		})
}
