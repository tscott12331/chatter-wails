//import { DebugLogger } from "@renderer/util/debug";
import { IAccessContextSuccess } from "@contexts/access-context";
import { TUser } from "@/App";
import { TApiResponse } from "./api";
import { FailedApiRequest, ServerError, Success, UnknownError } from "./api-response";
import { apiDeleteSubscriptions, apiGetSubscriptions, apiPostSubscriptions } from "./api-fetch";

//const dbLog = new DebugLogger();
export type TSubscriptionType = "channel.chat.message";

interface ICreateSubscriptionResponse {
    sessionId: string;
}

export const connectToChatroom = async (
                    userObj: TUser,
                    accessObj: IAccessContextSuccess,
                    brdId: string,
                    sesId: string,
                    setSubId?: React.Dispatch<React.SetStateAction<string | undefined>>,
                   ): Promise<TApiResponse<ICreateSubscriptionResponse>> =>
{
    return await createSubscription(userObj, accessObj, brdId, sesId, 'channel.chat.message', setSubId);
}

export const createSubscription = async (
                     userObj: TUser,
                     accessObj: IAccessContextSuccess,
                     brdId: string,
                     sesId: string,
                     subType: TSubscriptionType,
                     setSubId?: React.Dispatch<React.SetStateAction<string | undefined>>,
                     ): Promise<TApiResponse<ICreateSubscriptionResponse>> =>
{
    let retObj: TApiResponse<ICreateSubscriptionResponse> = UnknownError();

    try {
        const res = await apiPostSubscriptions(accessObj.access_token, {
            type: subType,
            version: "1",
            condition: {
                broadcaster_user_id: brdId.toString(),
                user_id: userObj.id.toString()
            },
            transport: {
                method: "websocket",
                session_id: sesId,
            },
        });

        if(res.status === 202) {
            const resData = await res.json();

            retObj = Success({ sessionId: resData.data[0].id});
        } else {

            retObj = FailedApiRequest();
        }
    } catch(err) {
        retObj = ServerError();
    } finally {
        setSubId?.(retObj.success ? retObj.data.sessionId : undefined);
        return retObj;
    }
}


interface IDeleteSubscriptionResponse {

}

export const deleteSubscription = async (subId: string, accessObj: IAccessContextSuccess): Promise<TApiResponse<IDeleteSubscriptionResponse>> => {
    let retObj: TApiResponse<IDeleteSubscriptionResponse> = UnknownError();

    try {
        const res = await apiDeleteSubscriptions(accessObj.access_token, {
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

export const deleteAllSubscriptions = async (
        access: IAccessContextSuccess,
    ): Promise<TApiResponse<IDeleteAllSubscriptionsResponse>> =>
{
    let retObj: TApiResponse<IDeleteAllSubscriptionsResponse> = UnknownError();

    try {
        let res = await getSubscriptions(access);
        if(res.success) {
            for(const subscription of res.data.subscriptions) {
                const delRes = await deleteSubscription(subscription.id, access);

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

export const getSubscriptions = async (
    access: IAccessContextSuccess
    ): Promise<TApiResponse<IGetSubscriptionsResponse>> =>
{
    let retObj: TApiResponse<IGetSubscriptionsResponse> = UnknownError();
    try {
        const res = await apiGetSubscriptions(access.access_token);

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
