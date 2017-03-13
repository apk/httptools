package httptools

import (
	"net/http"
	"fmt"
	"log"
	"io/ioutil"
	"os/exec"
)

type picreq struct {
	ch chan []byte
	size int
}

func to_s(x int) string {
	return fmt.Sprintf("%d", x)
}

func picserve(ch chan picreq) {
	for rq := range ch {
		cmd := exec.Command(
			"raspistill",
			"-t", "1000",
			"-w", to_s (9 * 4 * rq.size),
			"-h", to_s (9 * 3 * rq.size),
			"-mm", "matrix",
			"-o", "-")
		out, err := cmd.Output()
		if err != nil {
			log.Print("Exec:", err)
		}
		rq.ch <- out
	}
}

func PiCam(mux *http.ServeMux, path string) {

	ch := make(chan picreq)

	go picserve(ch)

	defhdlr := func (suf string, fac int) {
		mux.HandleFunc(path + suf, func(w http.ResponseWriter, r *http.Request) {
			_, err := ioutil.ReadAll(r.Body)
			if err == nil {
				rc := make(chan []byte)
				ch <- picreq{ch: rc, size: fac}
				s := <-rc
				w.Write([]byte(s))
			}
		})
	}

	defhdlr ("", 8 * 9);
	defhdlr ("/r", 4 * 9);
	defhdlr ("/s", 2 * 9);
	defhdlr ("/t", 15);
	defhdlr ("/u", 9);
	defhdlr ("/v", 5);
}
