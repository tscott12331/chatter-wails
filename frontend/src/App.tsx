import { useContext } from 'react'
import { GlobalContext } from './contexts/global-context'
import TabManager, { TAB_MANAGER_HEIGHT } from './components/tabs/tab-manager';
import Router from './router';
import { TooltipContext } from './contexts/tooltip-context';
import ToastManager from './components/util/toast/toast-manager';
import Tooltip from './components/util/tooltip';
import WindowControls, { WINDOW_CONTROLS_HEIGHT } from './components/window/window-controls';

export default function App() {
    const { user, toast } = useContext(GlobalContext);
    const { currentTooltip } = useContext(TooltipContext);


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
                <ToastManager toast={toast}/>
                {currentTooltip && <Tooltip data={currentTooltip}/>}
            </div>
        </div>
    )
}
