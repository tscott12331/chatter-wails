import { createContext, useCallback, useEffect, useMemo, useState } from "react";

import { Login } from "@wailsjs/chatter-wails/services/auth/authservice";
import { GetGlobalBadgeSets } from "@wailsjs/chatter-wails/services/badge/badgeservice";
import { ApiBadgeSet } from "@wailsjs/chatter-wails/internal/api/nativeApi";
import { DeleteAllSubscriptions } from "@wailsjs/chatter-wails/services/eventsub/eventsubservice";
import { AppUser } from "@wailsjs/chatter-wails/shared/types";
import { IToast } from "@/components/util/toast/toast";
import { assertDefined } from "@/util/assert";
import { CancelError } from "@wailsio/runtime";
import ToastManager from "@/components/util/toast/toast-manager";


interface IGlobalContext {
    user: AppUser|null;
    globalBadgeSets: ApiBadgeSet[];
    broadcastToast: (toast: IToast) => void;
    broadcastError: (err: any) => void;
    toast?: IToast;
    submitAccessToken?: (accessToken: string) => Promise<boolean>;
}

export const GlobalContext = createContext<IGlobalContext>({
    user: null,
    globalBadgeSets: [],
    broadcastToast: () => {},
    broadcastError: () => {},
});

export function GlobalContextProvider({
    children
}: { children: React.ReactNode }) {
    const [user, setUser] = useState<AppUser|null>(null);
    const [globalBadgeSets, setGlobalBadgeSets] = useState<ApiBadgeSet[]>([]);
    const [toast, setToast] = useState<IToast|undefined>(undefined);


    const submitAccessToken = useCallback(async (accessToken: string) => {
        try {
            const appUser = await Login(accessToken);

            setUser({...appUser} as AppUser);
            if(appUser) {
                initData();
            }

            localStorage.setItem('token', accessToken);

            return true;
        } catch(err) {
            broadcastError(err)
            return false;
        }
    }, []);

    const tryLoginSavedToken = () => {
        try {
            const token = localStorage.getItem('token')?.valueOf();
            if(!token) return;

            submitAccessToken(token);
        } catch(err) {
            broadcastError(err);
        }
    }

    const initData = () => {
        // TODO: make these emote events as well
        if(globalBadgeSets.length === 0) {
            GetGlobalBadgeSets()
                .then(gbs => {
                    assertDefined(gbs, "Global badge sets not found");
                    setGlobalBadgeSets(gbs);
                })
                .catch(broadcastError);
        }

        DeleteAllSubscriptions().catch(broadcastError);
    }

    const broadcastToast = useCallback((toast: IToast) => {
        setToast(toast);
    }, []);

    const broadcastError = useCallback((err: any) => {
        if(err instanceof CancelError) return;

        broadcastToast({
            message: `${err}`,
            type: 'error',
        });
    }, []);

    useEffect(() => {
        tryLoginSavedToken();
    }, []);

    const globalContextValue = useMemo(() => ({ 
            user,
            submitAccessToken,
            globalBadgeSets,
            broadcastToast,
            broadcastError,
            toast,
        }), [user, globalBadgeSets, toast]);

    return (
        <GlobalContext.Provider value={globalContextValue}>
            {children}
            <ToastManager toast={toast}/>
        </GlobalContext.Provider>
    )
}
