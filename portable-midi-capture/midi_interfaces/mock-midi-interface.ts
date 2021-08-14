import {MidiInterface} from './midi-interface';

export class MockMidiInterface implements MidiInterface {
    finishRecording = () => {}
    getActiveMidiInstruments = () => [];
    isRecordingMidi = () => false;
    record = () => {}
}
