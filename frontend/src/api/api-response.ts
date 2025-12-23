import { IApiFail, IApiSuccess } from "./api";
import { FailedApiRequestResponse, ServerErrorResponse, UnknownErrorResponse } from "./api-constants";

export const UnknownError = (): IApiFail => {
    return {...UnknownErrorResponse};
}

export const FailedApiRequest = (): IApiFail => {
    return {...FailedApiRequestResponse};
}

export const ServerError = (): IApiFail => {
    return {...ServerErrorResponse};
}

export const Success = <T>(data: T): IApiSuccess<T> => {
    return {
        success: true,
        data,
    }
}
