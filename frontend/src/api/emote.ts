import { AppEmoteSet } from "@wailsjs/chatter-wails/shared/types";


// provider -> section_id -> set
export interface IChatroomEmotes {
    providers: Map<string, Map<string, AppEmoteSet>>;
    _hash: number
}
