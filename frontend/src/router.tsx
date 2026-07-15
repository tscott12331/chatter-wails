import { Route, Routes } from "react-router-dom";
import Chatroom from "./pages/chatroom";
import HomePage from "./pages/home";
import NotFoundPage from "./pages/not-found";
import ToastManager from "./components/util/toast/toast-manager";
import { useContext } from "react";
import { GlobalContext } from "./contexts/global-context";

interface IRouterProps {
}

export default function Router({
}: IRouterProps) {
    const { toast } = useContext(GlobalContext);

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
        </>
    )
}
