import { TUser } from "@/App";
import { TApiResponse, TMessageResponse } from "./api";
import { apiPostMessages } from "./api-fetch";
import { FailedApiRequest, ServerError, Success, UnknownError } from "./api-response";


interface ISendMessageResponse  {
    message: TMessageResponse;
}

export const sendMessage = async (
        userObj: TUser,
        access_token: string,
        brdId: string,
        message: string,
        replyId?: string,
        ): Promise<TApiResponse<ISendMessageResponse>> =>
{
    let retObj: TApiResponse<ISendMessageResponse> = UnknownError();

    try {
        const res = await apiPostMessages(access_token, {
            sender_id: userObj.id,
            broadcaster_id: brdId,
            message,
            reply_parent_message_id: replyId,
        });
        if(res.ok) {
            const resData = await res.json();
            const messageResponse: TMessageResponse|undefined = resData.data?.at(0);
            if(messageResponse) {
                retObj = Success({ message: messageResponse });
            } else {
                retObj = FailedApiRequest();
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
