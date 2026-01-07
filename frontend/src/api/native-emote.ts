import { TUser } from "@/App";
import { IAccessContextSuccess } from "@contexts/access-context";
import { TApiResponse } from "./api";
import { FailedApiRequest, ServerError, Success, UnknownError } from "./api-response";
import { apiGetGlobalEmotes, apiGetUserEmotes } from "./api-fetch";

const scales: number[] = [1, 2, 4];

type TEmoteType = 'none'|'bitstier'|'follower'|'subscriptions'|'channelpoints'
                    |'rewards'|'hypetrain'|'prime'|'turbo'|'smilies'|'globals'
                    |'owl2019'|'twofactor'|'limitedtime';

export interface IUserEmote {
    id: string;
    name: string;
    emote_type: TEmoteType;
    emote_set_id: string;
    owner_id: string;
    format: ('animated'|'static')[];
    scale: number[];
    theme_mode: ('dark'|'light')[];
}

interface IGetUserEmotesResponse {
    emotes: IUserEmote[];
    template: string;
    pagination?: {
        cursor: string;
    }
}

export const getUserEmotes = async (
    userObj: TUser,
    accessObj: IAccessContextSuccess,
    setUserEmotes?: React.Dispatch<React.SetStateAction<IUserEmote[]>>,
): Promise<TApiResponse<IGetUserEmotesResponse>> => {
    let retObj: TApiResponse<IGetUserEmotesResponse> = UnknownError();

    try {
        const res = await apiGetUserEmotes(accessObj.access_token, {
            user_id: userObj.id,
        })

        if(res.ok) {
            const resData = await res.json();
            if('data' in resData && Array.isArray(resData.data)
              && 'template' in resData) {
                retObj = Success({
                        emotes: resData.data,
                        template: resData.template,
                        pagination: resData.pagination,
                    })
              } else {
                  retObj = FailedApiRequest();
              }
        } else {
            retObj = FailedApiRequest();
        }
    } catch(err) {
        console.error(err);
        retObj = ServerError();
    } finally {
        setUserEmotes?.(retObj.success ? retObj.data.emotes : []);
        return retObj;
    }
}


export interface IAppEmote {
	id: string
	name: string
	lightSrcSet: string
	darkSrcSet: string
}

export interface IGlobalEmote {
    id: string;
    name: string;
    images: {
        url_1x: string;
        url_2x: string;
        url_4x: string;
    };
    format: ('static'|'animated')[];
    scale: number[];
    theme_mode: ('light'|'dark')[];
}

interface IGetGlobalEmotesResponse {
    emotes: IGlobalEmote[];
    template: string;
}

export const getGlobalEmotes = async (
    accessObj: IAccessContextSuccess,
    setGlobalEmotes?: React.Dispatch<React.SetStateAction<IGlobalEmote[]>>,
): Promise<TApiResponse<IGetGlobalEmotesResponse>> => {
    let retObj: TApiResponse<IGetGlobalEmotesResponse> = UnknownError();

    try {
        const res = await apiGetGlobalEmotes(accessObj.access_token);

        if(res.ok) {
            const resData = await res.json();
            if('data' in resData && Array.isArray(resData.data)
              && 'template' in resData) {
                retObj = Success({
                        emotes: resData.data,
                        template: resData.template,
                    })
              } else {
                  retObj = FailedApiRequest();
              }
        } else {
            retObj = FailedApiRequest();
        }
    } catch(err) {
        console.error(err);
        retObj = ServerError();
    } finally {
        setGlobalEmotes?.(retObj.success ? retObj.data.emotes : []);
        return retObj;
    }
}


export const getEmoteSrcSet = (
        id: string,
        format: ('static'|'animated')[],
        color: 'light'|'dark' = "dark",
    ): string|undefined =>
{
    let srcSet = '';

    const prefFormat = format.includes('animated')
        ? 'animated'
        : format.includes('static')
        ? 'static'
        : undefined;
    if(!prefFormat) return undefined;

    const baseUrl = `https://static-cdn.jtvnw.net/emoticons/v2/${id}/${prefFormat}/${color}/`;
    scales.forEach((scale, i) => {
        const url = baseUrl.concat(`${scale}.0 ${scale}x`)
            .concat(i < scales.length - 1 ? ', ' : '');
        srcSet = srcSet.concat(url);
    })

    return srcSet;
}


