package gnord

import (
	"net/http"
	"net/http/cgi"
	"log"
	"os"
	"fmt"
	"flag"
	"path/filepath"
)

type GnordOpts struct {
	Path string
}

func GnordHandleFunc(opts *GnordOpts) func (w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		file := filepath.Join(opts.Path, filepath.FromSlash(path))
		ext := filepath.Ext(file)
		if ext == ".cgi" {
			// Hide cgi files from plain view
			http.NotFound(w, r)
			return
		}

		f, e := os.Lstat(file)
		if e == nil && (f.Mode() & os.ModeSymlink != 0) {
			s, e := os.Readlink(file)
			if e == nil {
				http.Redirect(w, r, s, http.StatusSeeOther)
				return
			}
		}

		if os.IsNotExist(e) {
			cginame := file + ".cgi"
			_, e = os.Stat(cginame)
			if (e == nil) {
				if *iphead != "" {
					ff := r.Header.Get(*iphead)
					if ff != "" {
						r.RemoteAddr = ff
					}
				}
				h := cgi.Handler{
					Path: cginame,
					Root: opts.Path,
				}
				h.ServeHTTP(w, r)
				return
			}
		}
		http.ServeFile(w, r, file)
	}
}
