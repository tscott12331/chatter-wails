import { TUser } from "@/App";
import { createContext, useEffect, useState } from "react";
import { Login } from '@wailsjs/go/services/AuthService';
import { GetGlobalBadgeSets } from "@wailsjs/go/services/BadgeService";
import { GetGlobalEmotes } from "@wailsjs/go/services/EmoteService";
import { IBadgeSet } from "@/api/badges";
import { IAppEmote } from "@/api/native-emote";
import { DeleteAllSubscriptions } from "@wailsjs/go/services/EventSubService";

interface IGlobalContext {
    user: TUser|null;
    globalBadgeSets: IBadgeSet[];
    globalEmotes: IAppEmote[];
    submitAccessToken?: (accessToken: string) => Promise<boolean>;
}

export const GlobalContext = createContext<IGlobalContext>({ user: null, globalBadgeSets: [], globalEmotes: [] });

export function GlobalContextProvider({
    children
}: { children: React.ReactNode }) {
    const [user, setUser] = useState<TUser|null>(null);
    const [globalBadgeSets, setGlobalBadgeSets] = useState<IBadgeSet[]>([]);
    const [globalEmotes, setGlobalEmotes] = useState<IAppEmote[]>([]);

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

        DeleteAllSubscriptions(accessToken);
    }

    useEffect(() => {
        tryLoginSavedToken();
    }, []);

    return (
        <GlobalContext.Provider value={{user, submitAccessToken, globalBadgeSets, globalEmotes}}>
            {children}
        </GlobalContext.Provider>
    )
}
