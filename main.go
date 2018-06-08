package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"unicode"

	"github.com/spf13/pflag"
)

const (
	envPrefix = "http-redirector"
)

var (
	toSite string
)

func toUpper(s string) string {
	res := new(bytes.Buffer)
	for _, r := range s {
		if unicode.IsPunct(r) || unicode.IsSymbol(r) {
			res.WriteRune('_')
		} else if unicode.IsLower(r) {
			res.WriteRune(unicode.ToUpper(r))
		} else {
			res.WriteRune(r)
		}
	}
	return res.String()
}

func envName(flagName string) string {
	return toUpper(envPrefix) + "_" + toUpper(flagName)

}

func setFlagsFromEnv() {
	pflag.VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			return
		}

		en := envName(f.Name)
		v, ok := os.LookupEnv(en)
		if !ok {
			return
		}
		pflag.Set(f.Name, v)
	})
}

func init() {
	pflag.StringVar(&toSite, "to-site", "", "site the incoming requests will be redirected to")
}

func main() {
	pflag.Parse()
	setFlagsFromEnv()

	if toSite == "" {
		log.Fatal("Usage: %s --to-site <site the requests will be redirected to>")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ff := r.Header.Get("X-Forwarded-For")
		if ff == "" {
			ff = "-"
		}

		ref := r.Referer()
		if ref == "" {
			ref = "-"
		}

		log.Printf("%s %s %s%s %s", r.RemoteAddr, ff, r.Host, r.RequestURI, ref)

		toURI := fmt.Sprintf("%s%s", toSite, r.RequestURI)
		http.Redirect(w, r, toURI, http.StatusMovedPermanently)
	})

	log.Printf("Redirect all requests to %s\n", toSite)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
