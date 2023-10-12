package main

import (
	"bytes"
	"fmt"
	"html/template"
)

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

func renderUI() string {
	type Data struct {
		State string
	}

	state := fmt.Sprintf("%v", enabled)

	var b bytes.Buffer
	webTemplate.Execute(&b, Data{State: state})

	return b.String()
}
