require('dotenv').config();

// Load the AWS SDK for Node.js
import AWS from 'aws-sdk';
import {NOTIFY_PHONE_NUMBER} from './constants';
// Set region
AWS.config.update({region: 'us-east-1'});

interface Publisher {
    publish(params: {
        Message: string;
        PhoneNumber: string;
    }): {
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
        return this.notify('Program is running!');
    }

    notifyRecording = () => {
        return this.notify('Program is recording!');
    }

    notifyStoppedRecording = () => {
        return this.notify('Program has stopped recording!');
    }
}

export class FakePublisher {
    notify = (message: string) => {
        console.log(`RECEIVED TEXT: ${message}`);
    }
}

const sns = new AWS.SNS({apiVersion: '2010-03-31'});

const notifier = new Notifier(sns);
