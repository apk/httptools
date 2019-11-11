package httptools

import (
	"net/http"
	"net/http/cgi"
	"os"
	"path/filepath"
)

type GnordOpts struct {
	Path string
	IpHeader string
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
		if e == nil {
			if (f.Mode() & os.ModeSymlink != 0) {
				s, e := os.Readlink(file)
				if e == nil {
					http.Redirect(w, r, s, http.StatusSeeOther)
					return
				}
			} else if (f.Mode() & os.ModeDir != 0) {
				idxname := file + "/index"
				cginame := idxname + ".cgi"

				_, e = os.Stat(cginame)
				if (e == nil) {
					// Mostly common code with stuff below.
					if opts.IpHeader != "" {
						ff := r.Header.Get(opts.IpHeader)
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

				_, e = os.Stat(idxname)
				if (e == nil) {
					http.ServeFile(w, r, idxname)
					return
				}
			}
		}

		if os.IsNotExist(e) {
			cginame := file + ".cgi"
			_, e = os.Stat(cginame)
			if (e == nil) {
				if opts.IpHeader != "" {
					ff := r.Header.Get(opts.IpHeader)
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

		// Serve by default handler.
		http.ServeFile(w, r, file)
	}
}
