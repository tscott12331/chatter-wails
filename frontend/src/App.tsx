import { useContext } from 'react'
import { GlobalContext } from './contexts/global-context'
import TabManager, { TAB_MANAGER_HEIGHT } from './components/tabs/tab-manager';
import Router from './router';
import WindowControls, { WINDOW_CONTROLS_HEIGHT } from './components/window/window-controls';

export default function App() {
    const { user } = useContext(GlobalContext);


    return (
        <div className="h-full overflow-hidden bg-chatter-bg text-chatter-text-primary">
            <WindowControls />
            {user &&
                <TabManager />
            }
            <div 
                className="sticky"
                style={{
                    height: `calc(100% - ${WINDOW_CONTROLS_HEIGHT + TAB_MANAGER_HEIGHT}px)`,
                    top: `${WINDOW_CONTROLS_HEIGHT+TAB_MANAGER_HEIGHT}px`,
                }}
            >
                <Router/>
            </div>
        </div>
    )
}
