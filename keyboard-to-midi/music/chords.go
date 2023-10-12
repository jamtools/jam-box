package music

import (
	"fmt"

	"github.com/mickmister/keyboard-to-midi/outputs"
	"gopkg.in/music-theory.v0/chord"
	"gopkg.in/music-theory.v0/note"
)

var midiNumberToNoteNameMap = map[int]string{
	0:  "C",
	1:  "C#",
	2:  "D",
	3:  "Eb",
	4:  "E",
	5:  "F",
	6:  "F#",
	7:  "G",
	8:  "G#",
	9:  "A",
	10: "Bb",
	11: "B",
}

var midiNumberToNoteClass = map[int]note.Class{
	0:  note.ClassNamed("C"),
	1:  note.ClassNamed("C#"),
	2:  note.ClassNamed("D"),
	3:  note.ClassNamed("Eb"),
	4:  note.ClassNamed("E"),
	5:  note.ClassNamed("F"),
	6:  note.ClassNamed("F#"),
	7:  note.ClassNamed("G"),
	8:  note.ClassNamed("G#"),
	9:  note.ClassNamed("A"),
	10: note.ClassNamed("Bb"),
	11: note.ClassNamed("B"),
}

const MinorThirdInterval = 3
const MajorThirdInterval = 4

func ConvertMidiNumbersToChord(midiNumbers []int) outputs.QueuedMidiChord {
	if len(midiNumbers) == 0 {
		return outputs.QueuedMidiChord{
			Name: "Invalid chord, no notes provided",
		}
	}

	root := midiNumbers[0]
	for _, num := range midiNumbers {
		if num < root {
			root = num
		}
	}

	rootRelative := mod(root)
	rootName := midiNumberToNoteNameMap[rootRelative]

	third := MajorThirdInterval

	notes := []outputs.Note{}

	for _, num := range midiNumbers {
		notes = append(notes, outputs.Note{
			MidiNumber: num,
		})

		thirdRelative := mod(num)
		diff := thirdRelative - rootRelative
		diffRelative := mod(diff)

		if diffRelative == MajorThirdInterval {
			third = MajorThirdInterval
		} else if diffRelative == MinorThirdInterval {
			third = MinorThirdInterval
		}
	}

	quality := "Major"
	if third == MinorThirdInterval {
		quality = "Minor"
	}

	name := fmt.Sprintf("%s %s", rootName, quality)

	return outputs.QueuedMidiChord{
		Chord: chord.Of(name),
		Name:  name,
		Notes: notes,
	}
}

func mod(num int) int {
	return num % 12
}
