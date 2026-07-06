import { createContext, useEffect, useState } from "react";

import { Login } from "@wailsjs/chatter-wails/services/auth/authservice";
import { GetGlobalBadgeSets } from "@wailsjs/chatter-wails/services/badge/badgeservice";
import { GetGlobalEmotes, GetUserEmotes } from "@wailsjs/chatter-wails/services/emote/emoteservice";
import { DeleteAllSubscriptions } from "@wailsjs/chatter-wails/services/eventsub/eventsubservice";
import { AppEmote } from "@wailsjs/chatter-wails/services/emote";
import { AppUser } from "@wailsjs/chatter-wails/services/auth";
import { ApiBadgeSet } from "@wailsjs/chatter-wails/internal/api";


interface IGlobalContext {
    user: AppUser|null;
    globalBadgeSets: ApiBadgeSet[];
    globalEmotes: AppEmote[];
    userEmotes: AppEmote[];
    submitAccessToken?: (accessToken: string) => Promise<boolean>;
}

export const GlobalContext = createContext<IGlobalContext>({ user: null, globalBadgeSets: [], globalEmotes: [], userEmotes: [] });

export function GlobalContextProvider({
    children
}: { children: React.ReactNode }) {
    const [user, setUser] = useState<AppUser|null>(null);
    const [globalBadgeSets, setGlobalBadgeSets] = useState<ApiBadgeSet[]>([]);
    const [globalEmotes, setGlobalEmotes] = useState<AppEmote[]>([]);
    const [userEmotes, setUserEmotes] = useState<AppEmote[]>([]);

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
                .then(gbs => setGlobalBadgeSets(gbs))
                .catch(e => console.error(e));
        }
        if(globalEmotes.length === 0) {
            GetGlobalEmotes(accessToken)
                .then(ge => setGlobalEmotes(ge))
                .catch(e => console.error(e));
        }
        if(userEmotes.length === 0) {
            GetUserEmotes(accessToken)
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
