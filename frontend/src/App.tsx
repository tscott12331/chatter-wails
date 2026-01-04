import { useContext, useState } from 'react'
import { AccessContext, type IAccessContextSuccess } from './contexts/access-context'
import TabManager from './components/tabs/tab-manager';
import Router from './router';
import { getUser } from './api/user-info';
import { getGlobalBadges, IBadgeSet } from './api/badges';
import { getGlobalEmotes, IGlobalEmote } from './api/native-emote';
import { ESCon, ESSubscription, TSubscriptionType, _deleteAllSubscriptions } from './api/eventsub';

import { services } from '@wailsjs/go/models'
import { Test } from '@wailsjs/go/main/app'
import { CreateSubscription } from "@wailsjs/go/services/EventSubService"

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
    access_token: string;
}

export default function App() {
    const [globalBadgeSets, setGlobalBadgeSets] = useState<IBadgeSet[]>([]);
    const [globalEmotes, setGlobalEmotes] = useState<IGlobalEmote[]>([]);
    const context = useContext(AccessContext);
    const [user, setUser] = useState<TUser|undefined>(undefined);

    const initData = async (access: IAccessContextSuccess) => {
        const userRes = await getUser(access.access_token);
        if(!userRes.success) {
            console.error(userRes.error);
            return;
        }
        setUser({...userRes.data.user, access_token: access.access_token});
        getGlobalBadges(access, setGlobalBadgeSets);

        getGlobalEmotes(access, setGlobalEmotes);

        await _deleteAllSubscriptions(access.access_token);
    }

    if(context && context.access && 'access_token' in context.access
      && user?.access_token !== context.access.access_token) {
        initData(context.access);
    }

    if(user) {
        getUser(user.access_token, 'yugi2x').then(ur => {
            if(!ur.success) return;

            let serviceUser = new services.User(user);

            CreateSubscription(serviceUser, {
                'broadcaster_user_id': ur.data.user.id,
                'user_id': user.id.toString(),
            },
            "channel.chat.message"
              ).then(r => {
                  const sub = new ESSubscription(r, "channel.chat.message");
                  ESCon.deleteSubscription(sub, user.access_token);
              }).catch(e => console.error(e));
        });
    }

    return (
        <></>
        // <div className="h-full overflow-hidden">
        //     {user &&
        //         <TabManager />
        //     }
        //     <div className="h-[calc(100%-36px)] flex flex-col">
        //         <Router
        //         user={user}
        //         globalBadgeSets={globalBadgeSets}
        //         globalEmotes={globalEmotes}
        //         />
        //     </div>
        // </div>
    )
}
