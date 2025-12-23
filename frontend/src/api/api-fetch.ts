import { paramsToString } from "@util/uri";
import { TQueryParamRecord } from "./api";
import { USERS_ENDPOINT, SUBSCRIPTIONS_ENDPOINT, VALIDATE_ENDPOINT, MESSAGES_ENDPOINT, BADGES_ENDPOINT, EMOTES_ENDPOINT } from "./api-constants"

//type TFetchMethod = "GET" | "POST" | "PUT" | "PATCH" | "DELETE";


export const apiDelete = async (
    url: string,
    headers: HeadersInit,
    params?: TQueryParamRecord,
    ): Promise<Response> =>
{
    const paramUrl = params ? url.concat(paramsToString(params)) : url;
    return await fetch(paramUrl,
                   {method: "DELETE", headers});
}

export const apiGet = async (
    url: string,
    headers: HeadersInit,
    params?: TQueryParamRecord,
    ): Promise<Response> =>
{
    const paramUrl = params ? url.concat(paramsToString(params)) : url;
    return await fetch(paramUrl,
                   {method: "GET", headers});
}

export const apiPost = async (
    url: string,
    headers: HeadersInit,
    body: object,
    params?: TQueryParamRecord,
    ): Promise<Response> =>
{
    const paramUrl = params ? url.concat(paramsToString(params)) : url;
    return await fetch(paramUrl, {method: "POST",
                       headers: {
                           ...headers,
                           'Content-Type': 'application/json',
                       },
                       body: JSON.stringify(body)});
}


export const apiGetValidate = async (
        access_token: string,
    ): Promise<Response> =>
{
    return await apiGet(VALIDATE_ENDPOINT, ApiValidateHeaders(access_token));
}


export const apiGetSubscriptions = async (
    access_token: string,
    params?: TQueryParamRecord) =>
{
    return await apiGet(
                        SUBSCRIPTIONS_ENDPOINT,
                        ApiSubscriptionHeaders(access_token),
                        params,
    );
}

export const apiPostSubscriptions = async (
    access_token: string,
    body: object,
    params?: TQueryParamRecord
    ) =>
{
    return await apiPost(
        SUBSCRIPTIONS_ENDPOINT,
        ApiSubscriptionHeaders(access_token),
        body,
        params,
    );
}

export const apiDeleteSubscriptions = async(
    access_token: string,
    params?: TQueryParamRecord
    ) =>
{
    return await apiDelete(
        SUBSCRIPTIONS_ENDPOINT,
        ApiSubscriptionHeaders(access_token),
        params,
    );
}


export const apiGetUsers = async(
    access_token: string,
    params?: TQueryParamRecord
    ) =>
{
    return await apiGet(
        USERS_ENDPOINT,
        ApiUsersHeaders(access_token),
        params,
    );
}

export const apiPostUsers = async (
    access_token: string,
    body: object,
    params?: TQueryParamRecord
    ) =>
{
    return await apiPost(
        USERS_ENDPOINT,
        ApiUsersHeaders(access_token),
        body,
        params,
    );
}


interface IPostMessagesBody {
    broadcaster_id: string;
    sender_id: string;
    message: string; // 500 max
    reply_parent_message_id?: string;
}

export const apiPostMessages = async (
    access_token: string,
    body: IPostMessagesBody,
    params?: TQueryParamRecord
    ) =>
{
    return await apiPost(
        MESSAGES_ENDPOINT,
        ApiMessagesHeaders(access_token),
        body,
        params,
    );
}


export const apiGetChannelBadges = async (
    access_token: string,
    params?: TQueryParamRecord
    ) =>
{
    return await apiGet(
        BADGES_ENDPOINT,
        ApiBadgesHeaders(access_token),
        params
    );
}

export const apiGetGlobalBadges = async (
    access_token: string,
    ) =>
{
    return await apiGet(
        `${BADGES_ENDPOINT}/global`,
        ApiBadgesHeaders(access_token),
    );
}


interface IApiGetUserEmotesParams extends TQueryParamRecord {
    user_id: string;
}

export const apiGetUserEmotes = async (
    access_token: string,
    params: IApiGetUserEmotesParams,
    ) =>
{
    return await apiGet(`${EMOTES_ENDPOINT}/user`, ApiEmotesHeaders(access_token), params);
}


export const apiGetGlobalEmotes = async (
    access_token: string,
    params?: TQueryParamRecord,
    ) =>
{
    return await apiGet(`${EMOTES_ENDPOINT}/global`, ApiEmotesHeaders(access_token), params);
}


const ApiHeaders = (access_token: string) => {
    return {
        'Authorization': `Bearer ${access_token}`,
        'Client-Id': import.meta.env.VITE_CLIENT_ID,
    }
}

const ApiSubscriptionHeaders = ApiHeaders;

const ApiUsersHeaders = ApiHeaders;

const ApiMessagesHeaders = ApiHeaders;

const ApiBadgesHeaders = ApiHeaders;

const ApiEmotesHeaders = ApiHeaders;

const ApiValidateHeaders = (access_token: string) => {
    return {
        "Authorization" : "OAuth " + access_token,
    }
}
