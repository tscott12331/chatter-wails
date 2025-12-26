import { IAccessContextSuccess } from "@contexts/access-context";
import { apiGetChannelBadges, apiGetGlobalBadges } from "./api-fetch";
import { TApiResponse } from "./api";
import { FailedApiRequest, ServerError, Success, UnknownError } from "./api-response";
import { IMessageBadge } from "@components/chat/chat-message";

export interface IBadgeSet {
    set_id: string;
    versions: {
        id: string;
        image_url_1x: string;
        image_url_2x: string;
        image_url_4x: string;
        title: string;
        description: string;
        click_action: string|null;
        click_url: string|null;
    }[];
}

interface IGetChannelBadgesResponse {
    sets: IBadgeSet[];
}

export const getChannelBadges = async (
        access_token: string,
        brdId: string,
        setBadgeSets?: React.Dispatch<React.SetStateAction<IBadgeSet[]>>,
    ): Promise<TApiResponse<IGetChannelBadgesResponse>> =>
{
    let retObj: TApiResponse<IGetChannelBadgesResponse> = UnknownError();

    try {
        const res = await apiGetChannelBadges(access_token, {
            broadcaster_id: brdId,
        });

        if(res.ok) {
            const resData = await res.json();
            if(resData?.data) {
                retObj = Success({ sets: resData.data });
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
        setBadgeSets?.(retObj.success ? retObj.data.sets : []);
        return retObj;
    }
}


interface IGetGlobalBadgesResponse {
    sets: IBadgeSet[];
}

export const getGlobalBadges = async (
        access: IAccessContextSuccess,
        setBadgeSets?: React.Dispatch<React.SetStateAction<IBadgeSet[]>>,
    ): Promise<TApiResponse<IGetGlobalBadgesResponse>> =>
{
    let retObj: TApiResponse<IGetGlobalBadgesResponse> = UnknownError();

    try {
        const res = await apiGetGlobalBadges(access.access_token);

        if(res.ok) {
            const resData = await res.json();
            if(resData?.data) {
                retObj = Success({ sets: resData.data });
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
        setBadgeSets?.(retObj.success ? retObj.data.sets : []);
        return retObj;
    }
}


interface IESBadge {
    set_id: string;
    id: string;
    info: string;
}


export const combineChannelGlobalSets = (cBadgeSets: IBadgeSet[], gBadgeSets: IBadgeSet[]) => {
    const subscriberSet = combineBadgeSets(
            cBadgeSets.find(s => s.set_id === 'subscriber'),
            gBadgeSets.find(s => s.set_id === 'subscriber'),
    );
    const combined = [...cBadgeSets, ...gBadgeSets];

    if(subscriberSet) {
        const subIndex = combined.findIndex(s => s.set_id === 'subscriber');
        combined[subIndex] = subscriberSet;
    }

    return combined;
}

export const combineBadgeSets = (
    set1: IBadgeSet|undefined,
    set2: IBadgeSet|undefined,
    ): IBadgeSet|undefined => {
    if(!set1 || !set2 || set1.set_id !== set2.set_id) return undefined;
    return {
        set_id: set1.set_id,
        versions: [...set1.versions, ...set2.versions],
    };
}


export const esBadgesToMessageBadges = (
    esBadges: IESBadge[],
    badgeSets: IBadgeSet[],
    ): IMessageBadge[] =>
{
    const messageBadges: IMessageBadge[] = [];

    for(const badge of esBadges) {
        let srcSet = '';
        const set = badgeSets.find(bs => bs.set_id === badge.set_id);
        if(!set) continue;

        const version = set.versions.find(v => v.id === badge.id);
        if(!version) continue;

        const urls = [version.image_url_1x, version.image_url_2x, version.image_url_4x];
        urls.forEach((url, i) => {
            srcSet = srcSet.concat(url)
                .concat(` ${Math.pow(2, i)}x${i < urls.length ? ', ' : ''}`);
        })

        messageBadges.push({ srcSet, info: badge.info, title: version.title });
    }

    return messageBadges;
}
