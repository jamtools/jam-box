import {MidiInterface} from './midi-interface';

export class MockMidiInterface implements MidiInterface {
    finishRecording = async () => {}
    getActiveMidiInstruments = async () => [];
    isRecordingMidi = async () => false;
    record = () => {}
}
