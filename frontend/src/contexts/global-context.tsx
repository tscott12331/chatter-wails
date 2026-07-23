import { createContext, useEffect, useMemo, useState } from "react";

import { Login } from "@wailsjs/chatter-wails/services/auth/authservice";
import { GetGlobalBadgeSets } from "@wailsjs/chatter-wails/services/badge/badgeservice";
import { GetGlobalEmotes, GetUserEmotes } from "@wailsjs/chatter-wails/services/emote/emoteservice";
import { DeleteAllSubscriptions } from "@wailsjs/chatter-wails/services/eventsub/eventsubservice";
import { AppEmote } from "@wailsjs/chatter-wails/shared/types";
import { AppUser } from "@wailsjs/chatter-wails/services/auth";
import { ApiBadgeSet } from "@wailsjs/chatter-wails/internal/api";
import { IToast } from "@/components/util/toast/toast";
import { assertDefined } from "@/util/assert";
import { CancelError } from "@wailsio/runtime";


interface IGlobalContext {
    user: AppUser|null;
    globalBadgeSets: ApiBadgeSet[];
    globalEmotes: AppEmote[];
    userEmotes: AppEmote[];
    broadcastToast: (toast: IToast) => void;
    broadcastError: (err: any) => void;
    toast?: IToast;
    submitAccessToken?: (accessToken: string) => Promise<boolean>;
}

export const GlobalContext = createContext<IGlobalContext>({
    user: null,
    globalBadgeSets: [],
    globalEmotes: [],
    userEmotes: [],
    broadcastToast: () => {},
    broadcastError: () => {},
});

export function GlobalContextProvider({
    children
}: { children: React.ReactNode }) {
    const [user, setUser] = useState<AppUser|null>(null);
    const [globalBadgeSets, setGlobalBadgeSets] = useState<ApiBadgeSet[]>([]);
    const [globalEmotes, setGlobalEmotes] = useState<AppEmote[]>([]);
    const [userEmotes, setUserEmotes] = useState<AppEmote[]>([]);
    const [toast, setToast] = useState<IToast|undefined>(undefined);


    const submitAccessToken = async (accessToken: string) => {
        try {
            const appUser = await Login(accessToken);

            setUser({...appUser} as AppUser);
            if(appUser) {
                initData(appUser.access_token);
            }

            localStorage.setItem('token', accessToken);

            return true;
        } catch(err) {
            broadcastError(err)
            return false;
        }
    }

    const tryLoginSavedToken = () => {
        try {
            const token = localStorage.getItem('token')?.valueOf();
            if(!token) return;

            submitAccessToken(token);
        } catch(err) {
            broadcastError(err);
        }
    }

    const initData = (accessToken: string) => {
        if(globalBadgeSets.length === 0) {
            GetGlobalBadgeSets(accessToken)
                .then(gbs => {
                    assertDefined(gbs, "Global badge sets not found");
                    setGlobalBadgeSets(gbs);
                })
                .catch(broadcastError);
        }
        if(globalEmotes.length === 0) {
            GetGlobalEmotes(accessToken)
                .then(ge => {
                    assertDefined(ge, "Global emotes not found");
                    setGlobalEmotes(ge);
                })
                .catch(broadcastError);
        }
        if(userEmotes.length === 0) {
            GetUserEmotes(accessToken)
                .then(ue => {
                    assertDefined(ue, "User emotes not found");
                    setUserEmotes(ue);
                })
                .catch(broadcastError);
        }

        DeleteAllSubscriptions(accessToken).catch(broadcastError);
    }

    const broadcastToast = (toast: IToast) => {
        setToast(toast);
    }

    const broadcastError = (err: any) => {
        if(err instanceof CancelError) return;

        broadcastToast({
            message: `${err}`,
            type: 'error',
        });
    }

    useEffect(() => {
        tryLoginSavedToken();
    }, []);

    const globalContextValue = useMemo(() => ({ 
            user,
            submitAccessToken,
            globalBadgeSets,
            globalEmotes,
            userEmotes,
            broadcastToast,
            broadcastError,
            toast,
        }), [user, globalBadgeSets, globalEmotes, userEmotes, toast]);

    return (
        <GlobalContext.Provider value={globalContextValue}>
            {children}
        </GlobalContext.Provider>
    )
}
