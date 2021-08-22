import {MidiDeviceInfo, MidiInterface} from './midi-interface';

import {ChildProcess, exec} from 'child_process';

import {SIGINT, SIGKILL} from 'constants';
import {AREC_CMD} from '../constants';
import {parseMidiDevices} from '../util';

const find = require('find-process');

const log = console.log;

export class RealMidiInterface implements MidiInterface {
    recordingProcess?: ChildProcess;

    finishRecording = async () => {
        const p = this.recordingProcess;
        this.recordingProcess = undefined;

        if (p) {
            let success = p.kill(SIGINT)

            // this.runCommand(['kill', p.pid as unknown as string])

            if (success) {
                log(`Killed (SIGINT) [${p.pid}]`);
            } else {
                const msg = `Failed to kill recording process [${p.pid}]`;
                log(msg);
                throw new Error(msg);
            }
        } else {
            throw new Error('Couldn\'t find recording process')
        }
    }

    getActiveMidiInstruments = async () => {
        const cmd = [AREC_CMD, '-l'];
        const [p, stdout] = await this.runCommandWithOutput(cmd);
        // console.log(stdout);
        const devices = parseMidiDevices(stdout);
        return devices;
    }

    isRecordingMidi = async () => {
        if (!this.recordingProcess) {
            return false;
        }

        const pids = await this.getPids(this.recordingProcess.pid);
        return Boolean(pids.length && pids[0]);
    }

    record = (device: MidiDeviceInfo, fileName: string) => {
        const cmd = [AREC_CMD, '-p', device.portNumber, fileName];
        this.recordingProcess = this.runCommand(cmd);
    }

    getPids = async (pid?: number): Promise<Array<number | undefined>> => {
        if (!pid) {
            return [];
        }

        const list = await find('pid', pid) as ChildProcess[];
        return list.map((p) => p.pid);
    }

    runCommandWithOutput = (cmd: string[]): Promise<[ChildProcess, string]> => new Promise((resolve, reject) => {
        const c = cmd.join(' ');
        console.log('Running command', c);
        const process = exec(c, (err, stdout, stderr) => {
            console.log(err);
            console.log(stdout);
            console.log(stderr);
            if (err) {
                reject(err);
                return;
            }

            resolve([process, stdout]);
        });
    });

    runCommand = (cmd: string[]): ChildProcess => {
        const c = cmd.join(' ');
        console.log('Running command', c);
        const process = exec(c, (err, stdout, stderr) => {
            console.log(err);
            console.log(stdout);
            console.log(stderr);
            if (err) {
                console.error('Error running command', err);
                return;
            }
        });

        return process;
    }
}
