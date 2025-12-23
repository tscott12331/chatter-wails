import styles from './home.module.css';

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
    console.log(authURL);
    


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
        <div className={styles.wrapper + ' flex-center'}>
        {
            typeof user === 'undefined' ?
            <div className={styles.connectWrapper + ' flex-column flex-align-center'}>
                <h1 className={styles.conToTwitch}>
                    <a 
                    className="underline text-blue-200 cursor-pointer"
                    onClick={() => BrowserOpenURL(authURL)}
                    >
                        Connect to twitch</a> to get your access token
                </h1>
                <form onSubmit={handleTokenSubmit}>
                    <input placeholder='access token' id='access_token' name='access_token' />
                    <button type='submit'>connect</button>
                </form>
            </div>
                :
                    <div
                className={styles.userInfoWrapper + ' fill flex-column flex-align-center'}
                >
                <img className='big-pfp' src={user.profile_image_url} />
                <h1>{user.display_name}</h1>
                </div>
        }
        </div>
    )
}
