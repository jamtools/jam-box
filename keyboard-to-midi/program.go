package main

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/mickmister/keyboard-to-midi/outputs"
)

type Keys []string

// var handlers = map[int]func(){
// 	0: nextChord,
// 	1: saveChord,
// 	2: removeChord,
// }

type Program struct {
	keymap           map[string]int
	midiDevices      map[string]*outputs.MidiDevice
	websocketClients []interface{}
}

func (p *Program) nextChord() {

}

func (p *Program) saveChord() {
	m := p.midiDevices["a"]
	keysHeldDown := m.GetHeldDownKeys()

}

func (p *Program) handleMidiEvent() {
	p.refreshUI()
}

func (p *Program) removeChord() {

}

func (p *Program) serialize() ([]byte, error) {
	return nil, nil
}

func (p *Program) deserialize([]byte) error {
	return nil
}

func (p *Program) refreshUI() {
	uiHTML := p.renderUI()
	p.sendWebsocketEvent(map[string]interface{}{
		"type": "render",
		"data": uiHTML,
	})
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

func (p *Program) renderUI() string {
	keyboards := p.getAvailableKeyboards()
	midiDevices := p.getAvailableMidiDevices()

	type Data struct {
		MidiDevices []*outputs.MidiDevice
	}

	var b bytes.Buffer
	midiDeviceHTMLTemplate.Execute(&b, Data{MidiDevices: []*outputs.MidiDevice{
		{
			Name: "Yamaha Digital Piano",
		},
	}})

	midiDeviceText := b.String()

	// keyboardText := "<select>"
	// for _, k := range keyboards {
	// 	keyboardText += fmt.Sprintf("<option>%s</option>", k.Name)
	// }
	// keyboardText += "</select>"

	return fmt.Sprintf(`
		<html>
			<body>
				%s
			</body>
		</html>
	`, midiDeviceText)
}

func (p *Program) sendWebsocketEvent(event map[string]interface{}) {
	for _, wsClient := range p.websocketClients {
		fmt.Printf("wsClient: %v\n", wsClient)
		// wsClient.sendEvent(event)
	}
}
