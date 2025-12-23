import { IAccessContextSuccess } from "@contexts/access-context";
import { TApiResponse } from "./api";
import { TUser } from "@/App";
import { FailedApiRequest, ServerError, Success, UnknownError } from "./api-response";
import { apiGetUsers } from "./api-fetch";

interface IGetUserResponse {
    user: TUser;
}

export const getUser = async (
        access: IAccessContextSuccess,
        channelName?: string,
    ): Promise<TApiResponse<IGetUserResponse>> =>
{
    let retObj: TApiResponse<IGetUserResponse> = UnknownError();

    const params = channelName ?
        {
            login: channelName,
        }
        : undefined;

    try {
        const res = await apiGetUsers(access.access_token, params)

        if(res.ok) {
            const resUser: TUser|undefined = (await res.json())?.data?.at(0);
            if(resUser) {
                retObj = Success({user: resUser});
            } else {
                retObj = {
                    success: false,
                    error: "Could not find user",
                };
            }
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
