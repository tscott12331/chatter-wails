import { IAccessContextSuccess } from "@contexts/access-context";
import { TApiResponse } from "./api";
import { FailedApiRequest, ServerError, Success, UnknownError } from "./api-response";
import { apiGetValidate } from "./api-fetch";

interface IValidateAccessResponse {
    client_id: string;
    login: string;
    scopes: string[];
    user_id: string;
    expires_in: number;
}

export const validateAccess = async (acc: IAccessContextSuccess): Promise<TApiResponse<IValidateAccessResponse>> => {
    let retObj: TApiResponse<IValidateAccessResponse> = UnknownError();

    try {
        const res = await apiGetValidate(acc.access_token);
        if(res.ok) {
            retObj = Success(await res.json());
        } else {
            retObj = FailedApiRequest();
        }
    } catch(err) {
        console.error(err);
        retObj = ServerError();
    } finally {
        return retObj;
    }
}
