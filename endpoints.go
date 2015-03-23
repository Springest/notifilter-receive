package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx/types"
)

type ResponseWriter interface {
	http.ResponseWriter

	// StatusCode returns the written status code, or 0 if none has been written yet.
	StatusCode() int
	// Written returns whether the header has been written yet.
	Written() bool
	// Size returns the size in bytes of the body written so far.
	Size() int
}

type appResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

// Don't need this yet because we get it for free:
func (w *appResponseWriter) Write(data []byte) (n int, err error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	size, err := w.ResponseWriter.Write(data)
	w.size += size
	return size, err
}

func (w *appResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *appResponseWriter) StatusCode() int {
	return w.statusCode
}

func (w *appResponseWriter) Written() bool {
	return w.statusCode != 0
}

func (w *appResponseWriter) Size() int {
	return w.size
}

func (w *appResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the ResponseWriter doesn't support the Hijacker interface")
	}
	return hijacker.Hijack()
}

func (w *appResponseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (w *appResponseWriter) Flush() {
	flusher, ok := w.ResponseWriter.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}

func trackTime(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s\n", name, elapsed)
}

func requestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// method=GET path=/jobs/833552.json format=json controller=jobs action=show status=200 duration=58.33 view=40.43 db=15.26
		start := time.Now()
		var arw appResponseWriter
		arw.ResponseWriter = w
		h.ServeHTTP(&arw, r)
		elapsed := time.Since(start) / time.Millisecond
		fmt.Printf("method=%s path=%s status=%d remote_ip=%s duration=%dms\n", r.Method, r.RequestURI, arw.StatusCode(), r.RemoteAddr, elapsed)
	})
}

func handleFavicon(w http.ResponseWriter, r *http.Request) {
	defer trackTime(time.Now(), "handleFavicon")
	http.Redirect(w, r, "/markie", http.StatusFound)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	defer trackTime(time.Now(), "handleIndex")

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal("ParseFiles", err)
	}

	notifiers := []Notifier{}
	err = db.Select(&notifiers, "SELECT * FROM notifiers")
	if err != nil {
		log.Fatal("db.Select notifiers ", err)
	}

	incoming := []Incoming{}
	err = db.Select(&incoming, "SELECT * FROM incoming ORDER BY id desc")
	if err != nil {
		log.Fatal("db.Select incoming ", err)
	}

	data := map[string]interface{}{
		"notifiers": notifiers,
		"incoming":  incoming,
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal("t.Execute", err)
	}
}

func handleNew(w http.ResponseWriter, r *http.Request) {
	defer trackTime(time.Now(), "handleNew")

	t, err := template.ParseFiles("templates/new.html")
	if err != nil {
		log.Fatal("ParseFiles", err)
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Fatal("t.Execute", err)
	}
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	defer trackTime(time.Now(), "handleCreateRule")

	err := r.ParseForm()
	if err != nil {
		log.Fatal("handleCreateRule", err)
	}
	fmt.Println("incoming parameters")
	fmt.Printf("%v\n", r.Form)

	notification_type := r.Form.Get("notification_type")
	class := r.Form.Get("class")
	template := r.Form.Get("template")
	rules := r.Form.Get("rules")

	_, err = db.NamedExec(`INSERT INTO notifiers (notification_type, class, template, rules) VALUES (:notification_type, :class, :template, :rules)`,
		map[string]interface{}{
			"notification_type": notification_type,
			"class":             class,
			"template":          template,
			"rules":             types.JsonText(rules),
		})

	if err != nil {
		log.Fatal("insert named query", err)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func handleCount(w http.ResponseWriter, r *http.Request) {
	defer trackTime(time.Now(), "handleCount")

	count := countRows()
	jmap := map[string]int{
		"status": 200,
		"count":  count,
	}

	output, err := json.MarshalIndent(jmap, "", "  ")
	if err != nil {
		log.Fatal("MarshalIndent", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(output)
}
