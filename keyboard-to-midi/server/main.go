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

type UIState struct {
	Title   string
	Body    string
	Buttons []string
}

var enabled = false

func toggleSomething() {
	enabled = !enabled
}

func main() {
	_, err := loadState()
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan string)

	var getState = func() State {
		return State{Enabled: enabled}
	}

	var getUIState = func() UIState {
		state := getState()
		return UIState{
			Title:   "Toggle the thing",
			Body:    fmt.Sprintf("Enabled: %v", state.Enabled),
			Buttons: []string{"Toggle", "Do something else"},
		}
	}

	var clicked = func(btn string) {
		defer saveState(getState())

		if btn == "0" {
			toggleSomething()
		}
	}

	s := &Server{}
	go s.StartServer(clicked, getUIState)

	for {
		buttonClicked := <-ch
		fmt.Println(buttonClicked)
	}
}

type Server struct {
	clicked    func(btn string)
	getUIState func() UIState
}

func (s *Server) StartServer(clicked func(btn string), getUIState func() UIState) {
	s.clicked = clicked
	s.getUIState = getUIState

	r := mux.NewRouter()
	r.HandleFunc("/", s.handleRoot)
	r.HandleFunc("/submit", s.handleSubmit)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	ui := s.renderUI()
	fmt.Fprint(w, ui)
}

func (s *Server) handleSubmit(w http.ResponseWriter, r *http.Request) {
	defer s.refreshUI()

	webButtonID := r.URL.Query().Get("id")
	s.clicked(webButtonID)

	ui := s.renderUI()
	fmt.Fprint(w, ui)
}

func (s *Server) refreshUI() {
	ui := s.renderUI()
	s.sendToAllWebsockets(ui)
}

func (s *Server) sendToAllWebsockets(ui string) {

}

const jsScript = `
<script>
	(() => {
		const handleFormSubmit = (e) => {
			e.preventDefault();
			console.log(e);

			fetch('/submit?id=' + e.target.id, {
				method: 'POST',
			}).then(r => r.text()).then(text => {
				document.body.innerHTML = text;
				addListeners();
			});
		};

		const addListeners = () => {
			Array.from(document.querySelectorAll('button')).forEach(b => b.addEventListener('click', handleFormSubmit));
		}

		if (!window.loaded) {
			window.addEventListener('load', addListeners);
			window.loaded = true;
		}
	})();
</script>
`

var webTemplate = template.Must(template.New("").Parse(fmt.Sprintf(`
			<h1>{{.State.Title}}</h1>
			<h1>{{.State.Body}}</h1>
			{{range $i, $btn := .State.Buttons}}
				<button id="{{$i}}" type='button'>{{$btn}}</button>
			{{end}}

			%s
`, jsScript)))

func (s *Server) renderUI() string {
	type Data struct {
		State UIState
	}

	uiState := s.getUIState()

	var b bytes.Buffer
	webTemplate.Execute(&b, Data{uiState})

	return b.String()
}

type State struct {
	Enabled bool `json:"enabled"`
}

const stateFileName = "state.json"

func saveState(state State) error {
	// state := State{Enabled: enabled}
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
