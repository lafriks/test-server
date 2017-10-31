/*
Copyright (c) 2017 Lauris Bukšis-Haberkorns <lauris@nix.lv>
Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// handleActualRequest handles simple cross-origin requests, actual request or redirects
func handleActualRequest(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	origin := r.Header.Get("Origin")

	if r.Method == "OPTIONS" {
		return
	}
	// Always set Vary, see https://github.com/rs/cors/issues/10
	headers.Add("Vary", "Origin")
	if origin == "" {
		return
	}
	headers.Set("Access-Control-Allow-Origin", origin)
	headers.Set("Access-Control-Allow-Credentials", "true")
}

// handlePreflight handles pre-flight CORS requests
func handlePreflight(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	origin := r.Header.Get("Origin")

	if r.Method != "OPTIONS" {
		return
	}
	// Always set Vary headers
	// see https://github.com/rs/cors/issues/10,
	//     https://github.com/rs/cors/commit/dbdca4d95feaa7511a46e6f1efb3b3aa505bc43f#commitcomment-12352001
	headers.Add("Vary", "Origin")
	headers.Add("Vary", "Access-Control-Request-Method")
	headers.Add("Vary", "Access-Control-Request-Headers")

	if origin == "" {
		return
	}

	reqMethod := r.Header.Get("Access-Control-Request-Method")
	reqHeaders := r.Header.Get("Access-Control-Request-Headers")
	headers.Set("Access-Control-Allow-Origin", origin)

	// Spec says: Since the list of methods can be unbounded, simply returning the method indicated
	// by Access-Control-Request-Method (if supported) can be enough
	headers.Set("Access-Control-Allow-Methods", strings.ToUpper(reqMethod))
	// Spec says: Since the list of headers can be unbounded, simply returning supported headers
	// from Access-Control-Request-Headers can be enough
	headers.Set("Access-Control-Allow-Headers", reqHeaders)
	headers.Set("Access-Control-Allow-Credentials", "true")
}

func CorsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
			handlePreflight(w, r)
			// Preflight requests are standalone and should stop the chain as some other
			// middleware may not handle OPTIONS requests correctly. One typical example
			// is authentication middleware ; OPTIONS requests won't carry authentication
			// headers (see #1)
			w.WriteHeader(http.StatusOK)
		} else {
			handleActualRequest(w, r)
			h.ServeHTTP(w, r)
		}
	})
}

func main() {
	var dir, bind string
	flag.StringVar(&dir, "dir", ".", "Serve directory")
	flag.StringVar(&bind, "bind", "", "Listen on IP address")

	flag.Parse()

	var port = flag.Arg(0)
	if port == "" {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n%s <port>\n", os.Args[0], os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
		return
	}
	if err := http.ListenAndServe(bind+":"+port, CorsHandler(http.FileServer(http.Dir(dir)))); err != nil {
		fmt.Fprintf(os.Stderr, "Error while starting server: %s\n", err.Error())
		os.Exit(1)
	}
}
