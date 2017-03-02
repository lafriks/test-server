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
)

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
	if err := http.ListenAndServe(bind+":"+port, http.FileServer(http.Dir(dir))); err != nil {
		fmt.Fprintf(os.Stderr, "Error while starting server: %s\n", err.Error())
		os.Exit(1)
	}
}
