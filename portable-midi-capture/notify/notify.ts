import {NOTIFY_PHONE_NUMBER} from '../constants';
import {getMidiFileURL, getPrivateIPAddress, getServerURLFromIPAddress} from '../util';

export type PublishParams = {
    Message: string;
    PhoneNumber: string;
};

export interface Publisher {
    publish(params: PublishParams): {
        promise: () => Promise<any>;
    }
}

export class Notifier {
    private publisher: Publisher;

    constructor(publisher: Publisher) {
        this.publisher = publisher;
    }

    notify = async (message: string) => {
        const params = {
            Message: message,
            PhoneNumber: NOTIFY_PHONE_NUMBER,
        };

        return this.publisher.publish(params).promise();
    }

    notifyRunning = () => {
        try {
            let fileMessage = 'No IP address found for serving files.';

            const ip = getPrivateIPAddress();
            if (ip) {
                const serverURL = getServerURLFromIPAddress(ip);
                fileMessage = `Access my files at ${serverURL}`;
            }

            return this.notify(`Program is running! ${fileMessage}`);
        } catch (e) {
            console.error(e);
            throw e;
        }
    }

    notifyStartedRecording = () => {
        return this.notify('Program is recording!');
    }

    notifyStoppedRecording = (fname: string) => {
        const ip = getPrivateIPAddress();
        let fileMessage = 'No IP found to serve files.'
        if (ip) {
            const fileURL = getMidiFileURL(ip, fname);
            fileMessage = `Download here: ${fileURL}`;
        }

        return this.notify(`Program has stopped recording! ${fileMessage}`);
    }

    notifyError = (s: string) => {
        s = s.substring(0, 160);
        return this.notify(s);
    }
}
