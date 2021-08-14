import {MidiInterface} from './midi-interface';

import {ChildProcess, exec} from 'child_process';
import util from 'util';
import {SIGINT} from 'constants';
import {getMidiDevices} from './parse_devices';

const find = require('find-process');

export class MockMidiInterface implements MidiInterface {
    finishRecording = () => { }
    getActiveMidiInstruments = () => [];
    isRecordingMidi = () => false;
    record = () => { }
}

export class RealMidiInterface {
    finishRecording = () => { }
    getActiveMidiInstruments = () => [];
    isRecordingMidi = () => false;
    record = () => { }

    getPorts = async () => {
        const cmd = [arec, '-l'];
        const [p, stdout] = await this.runCommandWithOutput(cmd);
        console.log(stdout);
        const devices = getMidiDevices(stdout);
        // const str = check_output([arec, '-l'], universal_newlines = True)
        // const res = [line.split()[0] for line in str.split('\n') if 'MIDI' in line]
        // return res
        return devices.slice(1); // get rid of midi through device
    }

    getPids = async (pid?: number): Promise<Array<number | undefined>> => {
        if (!pid) {
            return [];
        }

        const list = await find('pid', pid) as ChildProcess[];
        return list.map((p) => p.pid);
    }

    runCommandWithOutput = (cmd: string[]): Promise<[ChildProcess, string]> => new Promise((resolve, reject) => {
        const process = exec(cmd.join(' '), (err, stdout, stderr) => {
            // const process = exec(cmd.join(' '), {}, (err, stdout, stderr) => {
            if (err) {
                reject(err);
                return;
            }

            console.log(stdout);
            resolve([process, stdout]);
            // exec(command, function(error, stdout, stderr){ callback(stdout); });
        });
    });

    runCommand = (cmd: string[]): ChildProcess => {
        // runCommand = (cmd: string[]): Promise<[ChildProcess, string]> => new Promise((resolve, reject) => {
        const process = exec(cmd.join(' '), (err, stdout, stderr) => {
            // const process = exec(cmd.join(' '), {}, (err, stdout, stderr) => {
            if (err) {
                console.error('Error running command', err);
                return;
            }

            console.log(stdout);
            // exec(command, function(error, stdout, stderr){ callback(stdout); });
        });

        return process;
    }
}

const midi = new RealMidiInterface();

const recpath = '.';
// const recpath = '/home/pi/midi';
// const tz = pytz.timezone('Europe/Helsinki')
const arec = '/usr/bin/arecordmidi'

const log = console.log;
log('recmidi.sh started');


const Popen = (cmd: string[]) => {
    return {} as any;
}

const kill = (pid: any, signal: any) => {
}

const stdoutFlush = () => {

}

let numFiles = 0;

(async () => {
    let recordingProcess: ChildProcess | undefined = undefined;
    let stdout: string;
    while (true) {
        try {
            const pids = await midi.getPids(recordingProcess?.pid)
            const ports = await midi.getPorts()
            const timestr = 'itstime';
            // const timestr = datetime.now(tz).strftime('%Y-%m-%d_%H.%M.%S')

            if (ports?.length && pids?.length && pids[0]) {
                console.log(numFiles, 'Recording...');

                // We're still recording

            } else if (ports?.length) {
                numFiles++;
                const fname = `${recpath}/rec-${numFiles}.mid`;
                const cmd = [arec, '-p', ports[0].portNumber, fname];
                recordingProcess = midi.runCommand(cmd);
                // console.log(stdout);

                // We're starting recording

                log(timestr, `Started [${recordingProcess.pid}] ${cmd.join(' ')}`)
            } else if (pids?.length) {

                // We've unplugged the midi device

                try {
                    if (recordingProcess) {
                        recordingProcess.kill(SIGINT)
                        log(timestr, `Killed (SIGINT) [${recordingProcess.pid}]`);
                    } else {
                        throw new Error('Couldn\'t find recording process')
                    }
                } catch (e) {
                    throw new Error(`Error killing recording process: ` + e.message);
                }

                recordingProcess = undefined;
            } else {
                log('did nothing');
            }

            stdoutFlush();

            await new Promise(r => setTimeout(r, 1 * 1000));
            // time.sleep(10) // Sleep until retry
        } catch (e) {
            console.log('Error in loop', e);
        }
    }
})();
