//import { DebugLogger } from "@renderer/util/debug";
import { IAccessContextSuccess } from "@contexts/access-context";
import { TUser } from "@/App";
import { TApiResponse } from "./api";
import { FailedApiRequest, ServerError, Success, UnknownError } from "./api-response";
import { apiDeleteSubscriptions, apiGetSubscriptions, apiPostSubscriptions } from "./api-fetch";
import { IChatMessageFragment, IMessageReply, TChatMessage } from "@/components/chat/chat-message";
import { IESBadge } from "./badges";

//const dbLog = new DebugLogger();
export type TSubscriptionType = "channel.chat.message";

interface ICreateSubscriptionResponse {
    subId: string;
}

const _connectToChatroom = async (
    userObj: TUser,
    brdId: string,
    sesId: string,
): Promise<TApiResponse<ICreateSubscriptionResponse>> =>
{
    const condition: Map<string, string> = new Map([
        ['broadcaster_user_id', brdId.toString()],
        ['user_id', userObj.id.toString()]
    ]);
    return await _createSubscription(userObj, condition, sesId, 'channel.chat.message');
}

const _createSubscription = async (
    userObj: TUser,
    condition: Map<string, string>,
    sesId: string,
    subType: TSubscriptionType,
): Promise<TApiResponse<ICreateSubscriptionResponse>> =>
{
    let retObj: TApiResponse<ICreateSubscriptionResponse> = UnknownError();

    try {
        const res = await apiPostSubscriptions(userObj.access_token, {
            type: subType,
            version: "1",
            condition: Object.fromEntries(condition),
            transport: {
                method: "websocket",
                session_id: sesId,
            },
        });

        if(res.status === 202) {
            const resData = await res.json();

            retObj = Success({ subId: resData.data[0].id});
        } else {

            retObj = FailedApiRequest();
        }
    } catch(err) {
        retObj = ServerError();
    } finally {
        return retObj;
    }
}


interface IDeleteSubscriptionResponse {

}

const _deleteSubscription = async (subId: string, access_token: string): Promise<TApiResponse<IDeleteSubscriptionResponse>> => {
    let retObj: TApiResponse<IDeleteSubscriptionResponse> = UnknownError();

    try {
        const res = await apiDeleteSubscriptions(access_token, {
            id: subId
        })

        if(res.status === 204) {
            retObj = Success({});
        } else {
            retObj = FailedApiRequest();
        }
    } catch(err) {
        console.error(err);
        retObj = ServerError();
    } finally {
        return retObj;
    }
}

interface IDeleteAllSubscriptionsResponse {

}

export const _deleteAllSubscriptions = async (
    access_token: string,
): Promise<TApiResponse<IDeleteAllSubscriptionsResponse>> =>
{
    let retObj: TApiResponse<IDeleteAllSubscriptionsResponse> = UnknownError();

    try {
        let res = await _getSubscriptions(access_token);
        if(res.success) {
            for(const subscription of res.data.subscriptions) {
                const delRes = await _deleteSubscription(subscription.id, access_token);

                if(!delRes.success) {
                    return delRes;
                }
            }

            retObj = Success({});
        } else {
            retObj = res;
        }
    } catch(err) {
        console.error(err);
        retObj = ServerError();
    } finally {
        return retObj;
    }
}


interface ISubscription {
    id: string;
    status: string;
    type: string;
    version: string;
    condition: object;
    created_at: Date;
    transport: {
        method: string;
        callback: string;
    };
    cost: number;
}

interface IGetSubscriptionsResponse {
    subscriptions: ISubscription[];
}

export const _getSubscriptions = async (
    access_token: string
): Promise<TApiResponse<IGetSubscriptionsResponse>> =>
{
    let retObj: TApiResponse<IGetSubscriptionsResponse> = UnknownError();
    try {
        const res = await apiGetSubscriptions(access_token);

        if(res.ok) {
            const resData = await res.json();
            const subscriptions = resData.data;
            if(subscriptions && Array.isArray(subscriptions)) {
                retObj = Success({subscriptions});
            } else {
                retObj = {
                    success: false,
                    error: "API response mismatch",
                }
            }
        } else {
            retObj = FailedApiRequest();
        }
    } catch(err) {
        console.error(err);
        retObj = ServerError();
    } finally {
        return retObj;
    }
}



interface IESChatEvent {
    broadcaster_user_id: string;
    broadcaster_user_login: string;
    broadcaster_user_name: string;
    chatter_user_id: string;
    chatter_user_login: string;
    chatter_user_name: string;
    message_id: string;
    message: {
        text: string;
        fragments: IChatMessageFragment[];
    };
    color: string;
    badges: IESBadge[];
    message_type: string;
    cheer: {
        bits: number;
    }|null;
    reply: IMessageReply|null;
    channel_points_custom_reward_id: string|null;
    source_broadcaster_user_id: string|null;
    source_broadcaster_user_name: string|null;
    source_braodcaster_user_login: string|null;
    source_message_id: string|null;
    source_badges: IESBadge[]|null;
    is_source_only: boolean|null;
}

type TESEvent = IESChatEvent

interface IESWelcome {
    metadata: {
        message_id: string;
        message_type: string;
        message_timestamp: string;
    };
    payload: {
        session: {
            id: string;
            status: string;
            keepalive_timeout_seconds: number;
            reconnect_url: string;
            connected_at: string;
        }
    };
}

interface IESKeepalive {
    metadata: {
        message_id: string;
        message_type: string;
        message_timestamp: string;
    };
    payload: {};
}

export interface IESNotification {
    metadata: {
        message_id: string;
        message_type: string;
        message_timestamp: string;
        subscription_type: string;
        subscription_version: string;
    };
    payload: {
        subscription: {
            id: string;
            status: string;
            type: string;
            version: string;
            cost: number;
            condition: Record<string, string>;
            transport: {
                method: string;
                session_id: string;
            }
            created_at: string;
        }
        event: TESEvent;
    };
}

export type TESMessage = IESWelcome | IESKeepalive | IESNotification;


type TESSubscriptionEvent = 'message';
type TESSubscriptionHandler = (message: TESMessage) => void;

export class ESSubscription {
    private readonly handlers: Record<string, TESSubscriptionHandler[]> = {};

    addEventListener = (event: TESSubscriptionEvent, handler: TESSubscriptionHandler) => {
        if(event in this.handlers) {
            this.handlers[event].push(handler);
        } else {
            this.handlers[event] = [handler];
        }
    }

    removeEventListener = (event: TESSubscriptionEvent, handler: TESSubscriptionHandler) => {
        if(!(event in this.handlers)) return;

        const eventHandlers = this.handlers[event];
        const index = eventHandlers.indexOf(handler);
        if(index === -1) return;

        eventHandlers.splice(index, 1);
    }

    runEvent = (event: TESSubscriptionEvent, message: TESMessage) => {
        if(!(event in this.handlers)) return;

        for(const handler of this.handlers[event]) handler(message);
    }

    constructor(public subId: string, public subType: TSubscriptionType) {}

    destroy = () => {
        const keys = Object.keys(this.handlers);
        for(let i = keys.length - 1; i >= 0; i--) {
            const key = keys[i];
            delete this.handlers[key]
        }
    }
}

export class ESConnection {
    socket: WebSocket;
    sessionId: string|null = null;
    readonly subscriptions = {'channel.chat.message': new Set()} as Record<TSubscriptionType, Set<ESSubscription>>;


    createSubscription = async (
        userObj: TUser,
        condition: Map<string, string>,
        subType: TSubscriptionType,
    ): Promise<TApiResponse<{subscription: ESSubscription}>> => {
        if(!this.sessionId) return FailedApiRequest();

        const res = await _createSubscription(userObj, condition, this.sessionId, subType);
        if(!res.success) return FailedApiRequest();

        const subscription = new ESSubscription(res.data.subId, subType);

        this.subscriptions[subType].add(subscription);
        return Success({ subscription });
    }

    deleteSubscription = async (subscription: ESSubscription, access_token: string) => {
        subscription.destroy();
        this.subscriptions[subscription.subType].delete(subscription);
        return await _deleteSubscription(subscription.subId, access_token);
    }

    deleteAllSubscriptions = async (access_token: string) => {
        for(const [,subscriptionSet] of Object.entries(this.subscriptions)) {
            const subscriptionArr = [...subscriptionSet];
            for(let i = subscriptionArr.length - 1; i >= 0; i--) {
                const subscription = subscriptionArr[i];
                await this.deleteSubscription(subscription, access_token);
            }
        }
        
    }

    connectToChatroom = async (
        userObj: TUser,
        brdId: string,
    ): Promise<TApiResponse<{subscription: ESSubscription}>> => {
        const condition: Map<string, string> = new Map([
            ['broadcaster_user_id', brdId.toString()],
            ['user_id', userObj.id.toString()]
        ]);
        return this.createSubscription(userObj, condition, 'channel.chat.message');
    }

    handleNotificationMessage = (data: IESNotification) => {
        switch(data.metadata.subscription_type) {
            case 'channel.chat.message':
                this.subscriptions["channel.chat.message"].forEach(s => s.runEvent('message', data));
                break;
        }
    }

    handleEventsubMessage = (e: MessageEvent) => {
        const message: TESMessage = JSON.parse(e.data);
        const messageType = message.metadata.message_type;
        switch(messageType) {
            case 'session_welcome':
                this.sessionId = (message as IESWelcome).payload.session.id;
                break;
            case 'notification':
                this.handleNotificationMessage(message as IESNotification);
                break;
        }
    }

    constructor() {
        this.socket = new WebSocket("wss://eventsub.wss.twitch.tv/ws");

        this.socket.addEventListener('message', this.handleEventsubMessage)
    }

    destroy = (access_token: string) => {
        this.deleteAllSubscriptions(access_token);
    }
}



export const ESCon: ESConnection = new ESConnection();
