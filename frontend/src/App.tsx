import styles from './App.module.css';

import { useContext, useEffect, useState } from 'react'
import { AccessContext, type IAccessContextSuccess } from './contexts/access-context'
import TabManager from './components/tabs/tab-manager';
import Router from './router';
import { getUser } from './api/user-info';
import { deleteAllSubscriptions } from './api/eventsub';
import { getGlobalBadges, IBadgeSet } from './api/badges';
import { getGlobalEmotes, IGlobalEmote } from './api/native-emote';

export type TUser = {
    id: string;
    login: string;
    display_name: string;
    type: string;
    broadcaster_type: string;
    description: string;
    profile_image_url: string;
    offline_image_url: string;
    view_count: number;
    created_at: Date;
}

export default function App() {
    const [user, setUser] = useState<TUser>();
    const [globalBadgeSets, setGlobalBadgeSets] = useState<IBadgeSet[]>([]);
    const [globalEmotes, setGlobalEmotes] = useState<IGlobalEmote[]>([]);
    const context = useContext(AccessContext);

    const initUser = async (access: IAccessContextSuccess) => {
        const userRes = await getUser(access);
        if(!userRes.success) {
            console.error(userRes.error);
            return;
        }
        setUser(userRes.data.user);

        const delRes = await deleteAllSubscriptions(access);
        if(!delRes.success) {
            console.error(delRes.error);
        }
    }

    useEffect(() => {
        if(context && context.access && 'access_token' in context.access) {
            initUser(context.access);
            getGlobalBadges(context.access, setGlobalBadgeSets);
            getGlobalEmotes(context.access, setGlobalEmotes);
        }
    }, [context])

    return (
        <div className={styles.wrapper}>
            {user &&
                <TabManager />
            }
            <div className={styles.contentWrapper}>
                <Router
                user={user}
                globalBadgeSets={globalBadgeSets}
                globalEmotes={globalEmotes}
                />
            </div>
        </div>
    )
}
