import {MidiDeviceInfo} from './midi-interface';

export const getMidiDevices = (s: string): MidiDeviceInfo[] => {
    s = s.trim();

    let lines = s.split('\n');
    lines = lines.slice(2); // get rid of header and midi throughput

    lines = lines.filter(s2 => s2.length);
    lines = lines.map(l => l.trim());

    const rows = lines.map((line) => {
        const secondBegin = line.indexOf('  ');

        const portNumber = line.substring(0, secondBegin);

        line = line.substring(secondBegin);

        let second = 0;
        for (const [i, c] of Array.from(line).entries()) {
            if (c !== ' ') {
                second = i;
                break;
            }
        }

        line = line.substring(second);

        const thirdBegin = line.indexOf('  ');
        const clientName = line.substring(0, thirdBegin);
        line = line.substring(thirdBegin);

        let third = 0;
        for (const [i, c] of Array.from(line).entries()) {
            if (c !== ' ') {
                third = i;
                break;
            }
        }

        const portName = line.substring(third).trim();
        return {portNumber, clientName, portName};
    });

    return rows;
}

// const s = ` 14:0    Midi Through                     Midi Through Port-0
// 24:0    Launchkey Mini MK3               Launchkey Mini MK3 MIDI 1
// 24:1    Launchkey Mini MK3               Launchkey Mini MK3 MIDI 2`;

// console.log(getMidiDevices(s));
