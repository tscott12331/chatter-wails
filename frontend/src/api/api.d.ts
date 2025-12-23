export interface IApiFail {
    success: false;
    error: string;
}

export interface IApiSuccess<T> {
    success: true;
    data: T;
}

export type TApiResponse<T> = IApiFail|IApiSuccess<T>;

export type TQueryParamRecord = Record<string|number, string|number>;

export type TMessageResponse = {
    message_id: string;
    is_sent: true;
} | {
    message_id: string;
    is_sent: false;
    drop_reason?: {
        code: string;
        message: string;
    }
};

