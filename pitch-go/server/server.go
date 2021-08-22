package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/mickmister/pitch-go/view"
	"github.com/pkg/errors"
)

func handler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		err = errors.Wrap(err, "error parsing form")
		w.Write([]byte(err.Error()))
		return
	}

	var freq float32
	freqStr := r.Form.Get("freq")

	if len(freqStr) > 0 {
		freqFloat, err := strconv.ParseFloat(freqStr, 32)
		if err == nil {
			freq = float32(freqFloat)
			note := view.RenderFrequency(freq)
			SendMessage(string(note))
		}
	}

	s := getFrequencyForm(freq)
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprint(w, s)
}

func RunServer() {
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/", handler)

	fmt.Println("Serving from 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getFrequencyForm(freq float32) string {
	var noteName string
	noteName = view.RenderFrequency(freq)

	s := strconv.FormatFloat(float64(freq), 'f', 2, 32)

	return fmt.Sprintf(`
		<form method="post">
			<input name="freq" type="number" value="%s" step="0.01"></input>
			<input type="submit"></input>
			%s
		</form>
	`, s, noteName)
}
