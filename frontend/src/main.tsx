import './assets/main.css';

import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App'
import { HashRouter } from 'react-router-dom'
import { GlobalContextProvider } from '@/contexts/global-context'
import { TooltipContextProvider } from './contexts/tooltip-context';
import { TabContextProvider } from './contexts/tab-context';

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <HashRouter>
            <GlobalContextProvider>
                <TooltipContextProvider>
                    <TabContextProvider>
                        <App/>
                    </TabContextProvider>
                </TooltipContextProvider>
            </GlobalContextProvider>
        </HashRouter>
    </StrictMode>,
)
