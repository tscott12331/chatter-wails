import { AppEmoteSet } from "@wailsjs/chatter-wails/shared/types";


// provider -> section_id -> set
export type TChatroomEmotes = Map<string, Map<string, AppEmoteSet>>;
