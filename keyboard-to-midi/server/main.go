package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	ch := make(chan string)

	go StartServer(ch)

	for {
		c := <-ch
		fmt.Println(c)
	}
}

func StartServer(ch chan string) {
	_, err := loadState()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handleRoot)
	r.HandleFunc("/submit", handleSubmit(ch))

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	s := renderUI()
	fmt.Fprint(w, s)
}

func handleSubmit(ch chan string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer saveState()
		defer refreshUI()

		webButtonID := r.URL.Query().Get("id")
		if webButtonID == "1" {
			toggleSomething()
			ch <- webButtonID
		}

		s := renderUI()
		fmt.Fprint(w, s)
	}
}

var enabled = false

func toggleSomething() {
	enabled = !enabled
}

func refreshUI() {
	s := renderUI()
	sendToAllWebsockets(s)
}

func sendToAllWebsockets(s string) {

}

var midiDeviceHTMLTemplate = template.Must(template.New("").Parse(`
	<select>
	{{range $device := .MidiDevices}}
		<option
			value="{{$device.Name}}"
		>
			{{$device.Name}}
		</option>
	{{end}}
	</select>
`))

const jsScript = `
<script>
(() => {
	if (!window.loaded) {
		window.addEventListener('load', () => {
			document.querySelector('form').addEventListener('submit', handleFormSubmit);
		});
		window.loaded = true;
	}

		const handleFormSubmit = (e) => {
			e.preventDefault();
			console.log(e);

			fetch(e.target.action + '?id=1', {
				method: 'POST',
			}).then(r => r.text()).then(text => {
				document.body.innerHTML = text;
				document.querySelector('form').addEventListener('submit', handleFormSubmit);
			});
		};
	})();
</script>
`

var webTemplate = template.Must(template.New("").Parse(fmt.Sprintf(`
			<h1>{{.State}}</h1>
			<form action='/submit'>
				<input id="1" type='submit'>
			</form>

			%s
`, jsScript)))

func renderUI() string {
	type Data struct {
		State string
	}

	state := fmt.Sprintf("%v", enabled)

	var b bytes.Buffer
	webTemplate.Execute(&b, Data{State: state})

	return b.String()
}

type State struct {
	Enabled bool `json:"enabled"`
}

const stateFileName = "state.json"

func saveState() error {
	state := State{Enabled: enabled}
	b, err := json.Marshal(state)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(stateFileName, b, 0777)
}

func loadState() (state State, err error) {
	fmt.Println("Loading state")
	data, err := ioutil.ReadFile(stateFileName)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &state)

	enabled = state.Enabled

	fmt.Println(string(data))

	return
}
