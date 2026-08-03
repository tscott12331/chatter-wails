import { Route, Routes } from "react-router-dom";
import Chatroom from "./pages/chatroom";
import HomePage from "./pages/home";
import NotFoundPage from "./pages/not-found";
import SearchPage from "./pages/search";

interface IRouterProps {
}

export default function Router({
}: IRouterProps) {
    return (
        <>
        <Routes>
            <Route path='/' element={<HomePage />}/>
            <Route path='/chatroom/:channel' element={<Chatroom />} />
            <Route path='/search' element={<SearchPage />} />
            <Route path='/:any' element={<NotFoundPage />} />
        </Routes>
        </>
    )
}
