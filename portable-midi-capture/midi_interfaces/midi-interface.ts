export type MidiInstrument = {
    name: string;
};

export interface MidiInterface {
    getActiveMidiInstruments(): MidiInstrument[];
    record(): void;
    finishRecording(): void;
    isRecordingMidi(): boolean;
};
