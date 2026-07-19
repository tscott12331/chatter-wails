import { Route, Routes } from "react-router-dom";
import Chatroom from "./pages/chatroom";
import HomePage from "./pages/home";
import NotFoundPage from "./pages/not-found";
import ToastManager from "./components/util/toast/toast-manager";
import { useContext } from "react";
import { GlobalContext } from "./contexts/global-context";
import { TooltipContext } from "./contexts/tooltip-context";
import Tooltip from "./components/util/tooltip";

interface IRouterProps {
}

export default function Router({
}: IRouterProps) {
    const { toast } = useContext(GlobalContext);
    const { currentTooltip } = useContext(TooltipContext)

    return (
        <>
        <Routes>
            <Route path='/' element={<HomePage />}/>
            <Route path='/chatroom/:channel' element={
                <Chatroom />
                }
            />
            <Route path='/:any' element={<NotFoundPage />} />
        </Routes>
        <ToastManager toast={toast}/>
        {currentTooltip && <Tooltip data={currentTooltip}/>}
        </>
    )
}
