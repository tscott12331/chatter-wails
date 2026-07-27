import { useContext } from 'react';
import { GlobalContext } from '@/contexts/global-context';
import { Browser } from '@wailsio/runtime'

const scopes = [
    "user:read:chat",
    "user:write:chat",
    "user:read:emotes",
    "channel:moderate",
    "chat:read",
    "chat:edit",
];
const scopeParamVal = scopes.map(s => encodeURI(s)).join("+");

interface IHomePageProps {
}

export default function HomePage({
}: IHomePageProps) {
    const { user, submitAccessToken } = useContext(GlobalContext);
    const authURL = `https://id.twitch.tv/oauth2/authorize?response_type=token&client_id=${import.meta.env.VITE_CLIENT_ID}&redirect_uri=${import.meta.env.VITE_OAUTH_REDIRECT}&scope=${scopeParamVal}`;


    const handleTokenSubmit = (e: React.SubmitEvent<HTMLFormElement>) => {
        e.preventDefault();
        submitAccess(new FormData(e.currentTarget));
    }

    const submitAccess = (formData: FormData) => {
        const access_token = formData.get('access_token');
        if(!access_token) return;

        submitAccessToken?.(access_token.toString());
    }
    return (
        <div className='flex items-center justify-center h-full'>
        {
            user ?
            <div
            className='w-full h-full flex flex-col items-center gap-1 p-2.5'
            >
                <img className='big-pfp' src={user.profile_image_url} />
                <h1>{user.display_name}</h1>
            </div>
                :
            <div className={'flex flex-col items-center gap-4'}>
                <h1 className="w-1/2 text-center">
                    <a 
                    className="underline text-blue-200 cursor-pointer"
                    onClick={() => Browser.OpenURL(authURL)}
                    >
                        Connect to twitch</a> to get your access token
                </h1>
                <form onSubmit={handleTokenSubmit}>
                    <input className="m-2" placeholder='access token' id='access_token' name='access_token' />
                    <button className="m-2" type='submit'>connect</button>
                </form>
            </div>
        }
        </div>
    )
}
