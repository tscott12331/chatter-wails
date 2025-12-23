import { IApiFail } from "./api"

export const UnknownErrorResponse: IApiFail = {
    success: false,
    error: "Unknown error",
}

export const FailedApiRequestResponse: IApiFail = {
    success: false,
    error: "Failed request",
}

export const ServerErrorResponse: IApiFail = {
    success: false,
    error: "Server error",
}


export const VALIDATE_ENDPOINT = 'https://id.twitch.tv/oauth2/validate';

export const SUBSCRIPTIONS_ENDPOINT = 'https://api.twitch.tv/helix/eventsub/subscriptions';

export const USERS_ENDPOINT = 'https://api.twitch.tv/helix/users';

export const MESSAGES_ENDPOINT = 'https://api.twitch.tv/helix/chat/messages';

export const BADGES_ENDPOINT = 'https://api.twitch.tv/helix/chat/badges';

export const EMOTES_ENDPOINT = 'https://api.twitch.tv/helix/chat/emotes';
