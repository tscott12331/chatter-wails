import { Route, Routes } from "react-router-dom";
import Chatroom from "./pages/chatroom";
import HomePage from "./pages/home";
import NotFoundPage from "./pages/not-found";

interface IRouterProps {
}

export default function Router({
}: IRouterProps) {
    return (
        <>
        <Routes>
            <Route path='/' element={<HomePage />}/>
            <Route path='/chatroom/:channel' element={<Chatroom />}
            />
            <Route path='/:any' element={<NotFoundPage />} />
        </Routes>
        </>
    )
}
