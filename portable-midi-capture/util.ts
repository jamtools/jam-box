import {networkInterfaces} from 'os';

import {MidiDeviceInfo} from './midi_interfaces/midi-interface';

export const getTimeString = (): string => {
    const s = new Date().toISOString();
    return s.substring(0, s.length - 5).replaceAll(':', '-');
}

export const parseMidiDevices = (s: string): MidiDeviceInfo[] => {
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

export const getPrivateIPAddress = (): string | undefined => {
    const allNetworks = networkInterfaces();
    const results: {[name: string]: string[]} = {};

    for (const [name, nets] of Object.entries(allNetworks)) {
        if (!nets) {
            continue;
        }

        for (const net of nets) {
            // Skip over non-IPv4 and internal (i.e. 127.0.0.1) addresses
            if (net.family === 'IPv4' && !net.internal) {
                if (!results[name]) {
                    results[name] = [];
                }
                results[name].push(net.address);
            }
        }
    }

    if (results.en0) {
        return results.en0[0];
    }

    if (results.wlan0) {
        return results.wlan0[0];
    }

    return '';
}

export const getServerURLFromIPAddress = (ip: string): string => {
    return `http://${ip}:8000`;
}

export const getMidiFileURL = (ip: string, fname: string): string => {
    const serverURL = getServerURLFromIPAddress(ip);
    return `${serverURL}/data/${fname}`;
}
