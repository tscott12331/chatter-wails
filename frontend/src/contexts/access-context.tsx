import { validateAccess } from "@api/auth";
import { createContext, useEffect, useState } from "react";

const VALIDATE_INTERVAL: number = 86400000;


export interface IAccessContextSuccess {
    access_token: string;
    state?: string;
}

export interface IAccessContextError {
    error: string;
    error_description: string;
    state?: string;
}

interface IAccessContext {
    access: IAccessContextSuccess|IAccessContextError|null;
    setAccess: (newAccess: IAccessContextSuccess) => void;
}

export const AccessContext = createContext<IAccessContext|null>(null);

export function AccessContextProvider({
    children
}: { children: React.ReactNode }) {
    const [access, setAccess] = useState<IAccessContextSuccess|IAccessContextError|null>(null);

    const setAccessToken = (newAccess: IAccessContextSuccess) => {
        localStorage.setItem('token', JSON.stringify(newAccess));
        setAccess(newAccess);
    }

    const getAccessToken = () => {
        try {
            console.log('[getAccessCookie]: entered func');
            const cookieStr = localStorage.getItem('token')?.valueOf(); // FIX
            if(!cookieStr) return;

            const cookie: IAccessContextSuccess = JSON.parse(cookieStr);
            if(!cookie) return;
            console.log('[getAccessCookie]: got cookie', cookie);

            setAccess(cookie);
        } catch(err) {
            console.error(err);
        }
    }

    useEffect(() => {
        getAccessToken();
    }, []);

    useEffect(() => {
        if(access && 'access_token' in access) {
            validateAccess(access);
            const interval = setInterval(() => {
                validateAccess(access);
            }, VALIDATE_INTERVAL);

            return () => clearInterval(interval);
        }

        return;
    }, [access]);

    return (
        <AccessContext.Provider value={{access, setAccess: setAccessToken}}>
            {children}
        </AccessContext.Provider>
    )
}
