require('dotenv').config();

import AWS from 'aws-sdk';
import {Notifier} from './notify';

export const getNotifier = (): Notifier => {
    AWS.config.update({region: 'us-east-1'});
    const sns = new AWS.SNS({apiVersion: '2010-03-31'});
    return new Notifier(sns);
}
