import {Notifier, Publisher, PublishParams} from './notify';

export class FakePublisher implements Publisher {
    publish = (params: PublishParams) => {
        return {
            promise: async () => {
                console.log(`RECEIVED FAKE TEXT: ${params.Message}`);
            }
        }
    }
}

export const getNotifier = (): Notifier => {
    const publisher = new FakePublisher();
    return new Notifier(publisher);
}
