export interface IAppEmote {
	id: string
	name: string
	lightSrcSet: string
	darkSrcSet: string
    type: "global" | "user" | "channel"
    zeroWidth: boolean;
    emoteStack?: IAppEmote[];
}
