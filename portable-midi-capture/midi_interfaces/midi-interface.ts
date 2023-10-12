export type MidiDeviceInfo = {
    portNumber: string;
    clientName: string;
    portName: string;
}

export interface MidiInterface {
    getActiveMidiInstruments(): Promise<MidiDeviceInfo[]>;
    record(device: MidiDeviceInfo, fileName: string): void;
    finishRecording(): Promise<void>;
    isRecordingMidi(): Promise<boolean>;
};
