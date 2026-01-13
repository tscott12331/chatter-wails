import { useContext } from 'react'
import { GlobalContext } from './contexts/global-context'
import TabManager from './components/tabs/tab-manager';
import Router from './router';

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
    const { user } = useContext(GlobalContext);

    return (
        <div className="h-full overflow-hidden">
            {user &&
                <TabManager />
            }
            <div className="h-[calc(100%-36px)] flex flex-col">
                <Router/>
            </div>
        </div>
    )
}
