import {MockMidiInterface} from './midi_interfaces/mock-midi-interface';

const t = new MockMidiInterface();
t.record();

let i = 0;
for (const x of [1, 2, 3, 4, 5]) {
    i++;
    console.log(i);
}

t.finishRecording();
