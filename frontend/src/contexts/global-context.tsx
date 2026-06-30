import { TUser } from "@/App";
import { createContext, useEffect, useState } from "react";
import { IAppEmote } from "@/api/native-emote";

import { Login } from "@wailsjs/chatter-wails/services/auth/authservice";
import { GetGlobalBadgeSets } from "@wailsjs/chatter-wails/services/badge/badgeservice";
import { GetGlobalEmotes, GetUserEmotes } from "@wailsjs/chatter-wails/services/emote/emoteservice";
import { DeleteAllSubscriptions } from "@wailsjs/chatter-wails/services/eventsub/eventsubservice";


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

interface IGlobalContext {
    user: TUser|null;
    globalBadgeSets: IBadgeSet[];
    globalEmotes: IAppEmote[];
    userEmotes: IAppEmote[];
    submitAccessToken?: (accessToken: string) => Promise<boolean>;
}

export const GlobalContext = createContext<IGlobalContext>({ user: null, globalBadgeSets: [], globalEmotes: [], userEmotes: [] });

export function GlobalContextProvider({
    children
}: { children: React.ReactNode }) {
    const [user, setUser] = useState<TUser|null>(null);
    const [globalBadgeSets, setGlobalBadgeSets] = useState<IBadgeSet[]>([]);
    const [globalEmotes, setGlobalEmotes] = useState<IAppEmote[]>([]);
    const [userEmotes, setUserEmotes] = useState<IAppEmote[]>([]);

    const submitAccessToken = async (accessToken: string) => {
        try {
            const appUser = await Login(accessToken);

            setUser({...appUser});
            initData(appUser.access_token);


            localStorage.setItem('token', accessToken);

            return true;
        } catch(e) {
            console.error(e);
            return false;
        }
    }

    const tryLoginSavedToken = () => {
        try {
            const token = localStorage.getItem('token')?.valueOf();
            if(!token) return;

            submitAccessToken(token);
        } catch(err) {
            console.error(err);
        }
    }

    const initData = (accessToken: string) => {
        if(globalBadgeSets.length === 0) {
            GetGlobalBadgeSets(accessToken)
                // @ts-ignore
                .then(gbs => setGlobalBadgeSets(gbs))
                .catch(e => console.error(e));
        }
        if(globalEmotes.length === 0) {
            GetGlobalEmotes(accessToken)
                // @ts-ignore
                .then(ge => setGlobalEmotes(ge))
                .catch(e => console.error(e));
        }
        if(userEmotes.length === 0) {
            GetUserEmotes(accessToken)
                // @ts-ignore
                .then(ue => setUserEmotes(ue))
                .catch(e => console.error(e));
        }

        DeleteAllSubscriptions(accessToken);
    }

    useEffect(() => {
        tryLoginSavedToken();
    }, []);

    return (
        <GlobalContext.Provider value={{ user, submitAccessToken, globalBadgeSets, globalEmotes, userEmotes }}>
            {children}
        </GlobalContext.Provider>
    )
}
