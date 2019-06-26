package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/18f/cg-laboratories/networking/header-validator/libheader"
	"github.com/cloudfoundry-community/go-cfenv"
)

var (
	defaultPort          = 8080
	actualPort           string
	comparer             *libheader.HeaderComparer
	intender             *libheader.Intent
	headerRefFilename    string
	headerIntentFilename string
	defaultSeparator     string
)

func main() {
	log.Print("hello from boulder.")

	flag.StringVar(&headerRefFilename, "expectation-ref", "", `A JSON expectation reference in the format of map[string][]string{""} to use as the expected headers.`)
	flag.StringVar(&headerIntentFilename, "intent-ref", "", `A JSON intention in the format of a map[string][]string{""} to set headers that clients will expect.`)
	flag.StringVar(&defaultSeparator, "header-intent-separator", ",", `A single character used to separate the intent header fields for joining.`)
	flag.Parse()

	// load our reference file.
	headerRefFile, err := os.Open(headerRefFilename)
	if err != nil {
		log.Fatal(err)
	}
	comparer = libheader.NewComparer()
	if ok, err := comparer.Load(headerRefFile); !ok {
		log.Fatal(err)
	}

	// load our intent file.
	headerIntentFile, err := os.Open(headerIntentFilename)
	if err != nil {
		log.Fatal(err)
	}
	intender = libheader.NewIntent()
	if ok, err := intender.Load(headerIntentFile); !ok {
		log.Fatal(err)
	}

	appEnv, err := cfenv.Current()
	if err != nil {
		actualPort = fmt.Sprintf(":%d", defaultPort)
	} else {
		actualPort = fmt.Sprintf(":%d", appEnv.Port)
	}

	http.HandleFunc("/expect/", headerRepeater)
	http.HandleFunc("/expect/diff", headerDiff)
	http.HandleFunc("/intent/diff", headerDiffIntent)
	http.HandleFunc("/intent/", headerIntent)

	log.Printf("listening on %s", actualPort)
	log.Fatal(http.ListenAndServe(actualPort, nil))
}

func headerRepeater(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(r.Header)
}

func headerDiff(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	results := comparer.Compare(r.Header)
	// if they are identical, return teapot.
	if len(results) == 1 {
		w.WriteHeader(http.StatusTeapot)
	} else {
		enc.Encode(results)
	}
}

func headerIntent(w http.ResponseWriter, r *http.Request) {
	for key, val := range intender.Have {
		w.Header().Set(key, strings.Join(val, defaultSeparator))
	}
	log.Printf("sending headers %v", w.Header())
	json.NewEncoder(w).Encode(w.Header())
}

func headerDiffIntent(w http.ResponseWriter, r *http.Request) {
	// grab the internal hostname so we can hairpin.
	appHostname := url.URL{}
	if cfenv.IsRunningOnCF() {
		appEnv, err := cfenv.Current()
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		appHostname.Host = fmt.Sprintf("%s", appEnv.ApplicationURIs[0])
		appHostname.Scheme = "https"
	} else {
		appHostname.Host = fmt.Sprintf("%s%s", "localhost", actualPort)
		appHostname.Scheme = "http"
	}
	appHostname.Path = "/intent/"
	log.Printf("checking for intent validity at %s", appHostname.String())

	resp, err := http.DefaultClient.Get(appHostname.String())
	if err != nil {
		log.Print(err)
	}

	// instantiate a new temporary comparer so we can set our existing headers.
	localComparer := libheader.NewComparer()
	localComparer.Have = intender.Have

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	results := localComparer.Compare(resp.Header)
	// if they are identical, return teapot.
	if len(results) == 1 {
		w.WriteHeader(http.StatusTeapot)
	} else {
		enc.Encode(results)
	}
}
