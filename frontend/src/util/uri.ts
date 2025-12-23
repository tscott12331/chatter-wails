import { TQueryParamRecord } from "@renderer/api/api";

export type TParamReturn = Record<string, string|string[]|boolean>;

export function parseParams(paramString: string): TParamReturn {
    const pairs = paramString.split('&');
    const obj: TParamReturn = {};
    for(const pair of pairs) {
        const pairSplit = pair.split('=');
        const key = pairSplit.at(0);
        if(!key) continue;
        const value: string|boolean = pairSplit.at(1) ?? true

        obj[key] = value;
    }

    return obj;
}


export function paramsToString(params: TQueryParamRecord): string {
    let str = '?';
    for(const key in params) {
        str = str.concat(key).concat('=').concat(params[key].toString()).concat('&');
    }
    if(str.at(-1) === '&') str = str.slice(0, -1);
    return str;
}
