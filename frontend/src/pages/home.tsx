import type { TUser } from '../App';
import { useContext } from 'react';
import { AccessContext } from '@contexts/access-context';
import { BrowserOpenURL } from '@wailsjs/runtime/runtime'

interface IHomePageProps {
    user: TUser|undefined;
}

export default function HomePage({
    user,
}: IHomePageProps) {
    const context = useContext(AccessContext);
    const authURL = `https://id.twitch.tv/oauth2/authorize?response_type=token&client_id=${import.meta.env.VITE_CLIENT_ID}&redirect_uri=${import.meta.env.VITE_OAUTH_REDIRECT}&scope=user%3Aread%3Achat+user%3Awrite%3Achat+user%3Aread%3Aemotes`;


    const handleTokenSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        submitAccess(new FormData(e.currentTarget));
    }

    const submitAccess = (formData: FormData) => {
        if(!context) return;
        const access_token = formData.get('access_token');
        if(!access_token) return;

        const accessObj = { access_token: access_token.toString()};
        context.setAccess(accessObj);
    }
    return (
        <div className='flex items-center justify-center h-full'>
        {
            typeof user === 'undefined' ?
            <div className={'flex flex-col items-center gap-4'}>
                <h1 className="w-1/2 text-center">
                    <a 
                    className="underline text-blue-200 cursor-pointer"
                    onClick={() => BrowserOpenURL(authURL)}
                    >
                        Connect to twitch</a> to get your access token
                </h1>
                <form onSubmit={handleTokenSubmit}>
                    <input className="m-2" placeholder='access token' id='access_token' name='access_token' />
                    <button className="m-2" type='submit'>connect</button>
                </form>
            </div>
                :
                <div
                className='w-full h-full flex flex-col items-center gap-1 p-2.5'
                >
                <img className='big-pfp' src={user.profile_image_url} />
                <h1>{user.display_name}</h1>
                </div>
        }
        </div>
    )
}
