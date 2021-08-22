import {promises as fs} from 'fs';

import {RECORDING_PATH, SLEEPY_TIME} from './constants';

import {RealMidiInterface} from './midi_interfaces/real-midi-interface';
import {getTimeString} from './util';

import log from './logger';

import {Notifier} from './notify/notify';

import {getNotifier as awsNotify} from './notify/aws_notify';
import {getNotifier as fakeNotify} from './notify/fake_notify';
import {MidiInterface} from './midi_interfaces/midi-interface';
// import {getNotifier as soundtownNotify} from './soundtown';

// import {uploadFileToSoundTown} from './soundtown';

let notifier: Notifier;
if (process.argv[2] === 'sns') {
    log.info('Using SNS to send messages');
    notifier = awsNotify();
// } else if (process.argv[2] === 'soundtown') {
        // log.info('Using SoundTown to send messages');
        // notifier = soundtownNotify();
} else {
    log.info('Using console to send messages');
    notifier = fakeNotify();
}

let numFiles = 0;
const midi: MidiInterface = new RealMidiInterface();

(async () => {
    log.info('Loop started');
    notifier.notifyRunning();

    // await fs.mkdir(RECORDING_PATH, {recursive: true});

    let fname = '';
    let fpath = '';
    while (true) {
        try {
            const isRecording = await midi.isRecordingMidi();
            const devices = await midi.getActiveMidiInstruments();

            if (devices.length && isRecording) {
                log.debug('Still recording', numFiles);

            } else if (devices.length) {
                log.info('New device was plugged in');

                numFiles++;

                const timestr = getTimeString();
                fname = `${timestr}-${numFiles}.mid`;
                fpath = `${RECORDING_PATH}/${fname}`;
                midi.record(devices[0], fpath);
                // midi.record(devices[1], fpath + '2');

                log.info('Started recording');
                notifier.notifyStartedRecording();
            } else if (isRecording) {
                log.info('Device was unplugged. Stopping the recording.');

                try {
                    midi.finishRecording();
                } catch (e) {
                    throw new Error(`Error killing recording process: ` + e.message);
                }

                log.info('Stopped recording');
                notifier.notifyStoppedRecording(fname);

                // uploadFileToSoundTown(fpath)
            } else {
                log.debug('No devices plugged in. Waiting...');
            }

            await new Promise(r => setTimeout(r, SLEEPY_TIME));
        } catch (e) {
            log.error('Error in loop', e);
            notifier.notifyError(e.message);
        }
    }
})();
