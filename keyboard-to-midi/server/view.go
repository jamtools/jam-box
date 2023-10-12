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

var webTemplate = template.Must(template.New("").Parse(`
	<html>
		<body>
			<script>
				const handleFormSubmit = (e) => {
					console.log(e);
				}
			</script>

			<h1>{{.State}}</h1>
			<form action='/submit' onsubmit="handleFormSubmit">
				<input id="1" type='submit'>
			</form>
		</body>
	</html>
`))

func renderUI() string {
	type Data struct {
		State string
	}

	state := fmt.Sprintf("%v", enabled)

	var b bytes.Buffer
	webTemplate.Execute(&b, Data{State: state})

	return b.String()
}
