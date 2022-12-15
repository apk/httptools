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

// https://golang.org/src/net/http/fs.go
// See notes around about proper URL.Path handling.

// https://cs.opensource.google/go/go/+/refs/tags/go1.17.3:src/net/http/cgi/host.go;l=57

func GnordHandleFunc(opts *GnordOpts) func (w http.ResponseWriter, r *http.Request) {

	do404 := func (w http.ResponseWriter, r *http.Request) {
		file := filepath.Join(opts.Path, "404")
		http.ServeFile(w, r, file) // TODO: This returns a 200, not a 404!
	}

	docgi := func (cginame string, w http.ResponseWriter, r *http.Request) {
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
	}

	return func (w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		fp := filepath.FromSlash(path)
		file := filepath.Join(opts.Path, fp)
		ext := filepath.Ext(file)
		if ext == ".cgi" {
			// Hide cgi files from plain view
			do404 (w,r)
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
					docgi (cginame, w, r)
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
			_, e := os.Stat(cginame)
			if (e == nil) {
				docgi (cginame, w, r)
				return
			}
			//if (os.IsNotExist(e)) {
			//	do404(w,r)
			//	return
			//}

			for {
				fd := filepath.Dir(fp)
				if (fd == fp) { break }
				fp = fd

				cginame := filepath.Join(opts.Path, fp) + "/index.cgi"
				_, e := os.Stat(cginame)
				if (e == nil) {
					docgi (cginame, w, r)
					return
				}
			}
		}

		// Serve by default handler.
		http.ServeFile(w, r, file)
	}
}
