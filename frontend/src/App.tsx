import { useContext } from 'react'
import { GlobalContext } from './contexts/global-context'
import TabManager from './components/tabs/tab-manager';
import Router from './router';

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
