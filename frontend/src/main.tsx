import './assets/main.css';

import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App'
import { HashRouter } from 'react-router-dom'
import { AccessContextProvider } from './contexts/access-context';

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <HashRouter>
            <AccessContextProvider>
                <App/>
            </AccessContextProvider>
        </HashRouter>
    </StrictMode>,
)
