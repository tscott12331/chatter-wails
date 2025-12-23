import { Route, Routes } from "react-router-dom";
import Chatroom from "./pages/chatroom";
import HomePage from "./pages/home";
import type { TUser } from "./App";
import { IBadgeSet } from "./api/badges";
import NotFoundPage from "./pages/not-found";
import { IGlobalEmote } from "./api/native-emote";

interface IRouterProps {
    user: TUser|undefined;
    globalBadgeSets: IBadgeSet[];
    globalEmotes: IGlobalEmote[];
}

export default function Router({
    user,
    globalBadgeSets,
    globalEmotes,
}: IRouterProps) {
    return (
        <Routes>
            <Route path='/' element={<HomePage user={user}/>}/>
            {user &&
            <Route path='/chatroom/:channel' element={
                <Chatroom
                user={user}
                globalBadgeSets={globalBadgeSets}
                globalEmotes={globalEmotes}
                />
                }
            />
            }
            <Route path='/:any' element={<NotFoundPage />} />
        </Routes>
    )
}
