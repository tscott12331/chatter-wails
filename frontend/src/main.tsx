import './assets/main.css';

import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App'
import { HashRouter } from 'react-router-dom'
import { GlobalContextProvider } from '@/contexts/global-context'

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <HashRouter>
            <GlobalContextProvider>
                <App/>
            </GlobalContextProvider>
        </HashRouter>
    </StrictMode>,
)
