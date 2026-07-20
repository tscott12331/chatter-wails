import { useContext } from 'react'
import { GlobalContext } from './contexts/global-context'
import TabManager from './components/tabs/tab-manager';
import Router from './router';
import { TooltipContext } from './contexts/tooltip-context';
import ToastManager from './components/util/toast/toast-manager';
import Tooltip from './components/util/tooltip';
import WindowControls from './components/window/window-controls';

export default function App() {
    const { user, toast } = useContext(GlobalContext);
    const { currentTooltip } = useContext(TooltipContext);


    return (
        <div className="h-full overflow-hidden">
            <WindowControls />
            {user &&
                <TabManager />
            }
            <div className="h-[calc(100%-44px)] flex flex-col">
                <Router/>
                <ToastManager toast={toast}/>
                {currentTooltip && <Tooltip data={currentTooltip}/>}
            </div>
        </div>
    )
}
