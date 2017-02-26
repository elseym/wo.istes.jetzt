package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/elseym/go-tzlib"
	"github.com/elseym/go-tzlib/exporters"
	"github.com/elseym/go-tzlib/importers"
	"github.com/elseym/wo.istes.jetzt/responders"
)

var cfg = struct {
	tzlib    string
	from     string
	webroot  string
	endpoint string
	bind     string
}{}

// init loads the configuration
func init() {
	flag.StringVar(&cfg.tzlib, "tzlib", "./tzlib.json", "Use database file <tzlib>")
	flag.StringVar(&cfg.from, "from", "", "Import HTML from <from> and write database to <zlib>")
	flag.StringVar(&cfg.webroot, "webroot", "", "Serve the <webroot> directory at '/'")
	flag.StringVar(&cfg.endpoint, "endpoint", "/dann", "Serve the API at <endpoint>")
	flag.StringVar(&cfg.bind, "bind", "localhost:1620", "Listen on <bind> for connections. ")

	flag.Parse()
}

// main starts wo.istes.jetzt
func main() {
	if cfg.from != "" {
		update(cfg.from, cfg.tzlib)
	}

	log.Print(serve(cfg.bind, cfg.tzlib, cfg.endpoint, cfg.webroot))
}

// serve registers handlers and waits for connections
func serve(bind, tzlibfile, endpoint, webroot string) error {
	handleAPI(endpoint, tzlibfile)
	// handleWebroot(webroot)

	log.Printf("listening on 'http://%s'...", bind)
	return http.ListenAndServe(bind, nil)
}

// handleAPI sets up the tzlib api server
func handleAPI(at, jsonfile string) {
	tl := loadlib(jsonfile)

	if !strings.HasSuffix(at, "/") {
		at += "/"
	}
	handler := responders.DannResponder(tl)

	http.Handle(at, http.StripPrefix(at, handler))
}

// handleWebroot sets up the static file server
func handleWebroot(webroot string) {
	defer func() {
		if r := recover(); r != nil {
			log.Fatal("when using --webroot, --endpoint must not equal '/': ", r)
		}
	}()

	_, err := os.Stat(webroot)
	if err != nil && webroot != "" {
		log.Fatal("webroot '", webroot, "' inaccessible: ", err)
	}

	if webroot != "" {
		handler := responders.PublicResponder(webroot)
		http.Handle("/", handler)
	}
}

// loadlib is a convenient tzlib loader
func loadlib(jsonfile string) *tzlib.Tzlib {
	tl, err := tzlib.Import(importers.NewJSONFile(jsonfile))
	if err != nil {
		log.Fatal("could not read timezones from'", jsonfile, "': ", err)
	}

	return tl
}

// update takes parsable html from file <from> or url <from>,
// then creates a new tzlib and writes it to file <to>
func update(from, to string) {
	var importer tzlib.Importer
	var err error
	var l *tzlib.Tzlib

	_, err = os.Stat(from)
	if err == nil {
		importer = importers.NewTimeIsFromFile(from)
	} else {
		importer = importers.NewTimeIsFromURL(from)
	}

	l, err = tzlib.Import(importer)
	if err != nil {
		log.Fatalf("error importing '%s': %s", from, err.Error())
	}

	err = l.Export(exporters.NewJSONFile(to))
	if err != nil {
		log.Fatalf("error exporting '%s': %s", to, err.Error())
	}

	log.Printf("update successful: got %d timezones, data expires %s", len(l.Timezones), l.Expires.Format(time.RFC822Z))
}
