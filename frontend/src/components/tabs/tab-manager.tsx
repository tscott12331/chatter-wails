import { useContext, useEffect, useRef, useState } from 'react';
import Tab from './tab';
import HomeIcon from '../svg/home-icon';
import PlusIcon from '../svg/plus-icon';
import { useNavigate } from 'react-router-dom';
import { Events } from "@wailsio/runtime";
import { SharedChatParticipant } from '@wailsjs/chatter-wails/services/eventsub';
import SearchIcon from '../svg/search-icon';
import { TabContext } from '@/contexts/tab-context';
export type TTab = {
    tabRoute: string;
    tabName: string;
}

interface ITabManagerProps {

}

export default function TabManager({

}: ITabManagerProps) {
    const tabsRef = useRef<Record<string, {current: HTMLDivElement|null}>>({})

    const [isAddingTab, setIsAddingTab] = useState<boolean>(false);
    const [newTabText, setNewTabText] = useState<string>("");
    const [currentTabRoute, setCurrentTabRoute] = useState<string>('/');

    const { home, tabs, addTab, removeTab, editTab, rotateTabs } = useContext(TabContext)

    const navigate = useNavigate();

    const handleTabSelect = (tab: TTab) => {
        setCurrentTabRoute(tab.tabRoute);

        navigate(tab.tabRoute);
    }

    const handleAddTab = (tab: TTab) => {
        if(tabs.find(t => t.tabRoute.toLowerCase() === tab.tabRoute.toLowerCase())) {
            return handleTabSelect(tab)
        }

        addTab(tab);
    }

    const handleTabRemove = (tab: TTab) => {
        if(tab.tabRoute.toLowerCase() === currentTabRoute.toLowerCase()) {
            setCurrentTabRoute(home.tabRoute);
            navigate(home.tabRoute);
        }

        removeTab(tab);
        delete tabsRef.current[tab.tabRoute];
    }

    const addParticipantsToTabName = (channel: string, participants: Record<string, SharedChatParticipant|null|undefined>) => {
        const eventRoute = `/chatroom/${channel.toLowerCase()}`;
        const tabToChangeIndex = tabs.findIndex(t => t.tabRoute.toLowerCase() === eventRoute);
        editTab(tabToChangeIndex, (tab) => ({
            ...tab,
            tabName: Object.values(participants)
                        .map(p => p?.name ?? "unknown")
                        .reduce((p, c) => `${p}, ${c}`),
        }));
    }

    const handleSharedChatBegin = (event: Events.WailsEvent<"common:shared-chat-begin">) => {
        if(!event.data.participant) return tabs;
        addParticipantsToTabName(event.data.channel, event.data.participant);
    }

    const handleSharedChatUpdate = (event: Events.WailsEvent<"common:shared-chat-update">) => {
        if(!event.data.participant) return tabs;
        addParticipantsToTabName(event.data.channel, event.data.participant);
    }

    const handleSharedChatEnd = (event: Events.WailsEvent<"common:shared-chat-end">) => {
        const eventRoute = `/chatroom/${event.data.channel.toLowerCase()}`;
        const tabToChangeIndex = tabs.findIndex(t => t.tabRoute.toLowerCase() === eventRoute);
        editTab(tabToChangeIndex, (tab) => ({
            ...tab,
            tabName: eventRoute.split('/chatroom/')[1],
        }));
    }

    const handleAddTabKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if(e.key === 'Enter') {
            const tabName = newTabText.trim();
            if(!tabName.includes(' ') && tabName.length >= 3
               && tabName.length <= 25) {
                   const tabRoute = `/chatroom/${tabName.toLowerCase()}`;
                   const newTab: TTab = {
                        tabRoute,
                        tabName: tabName,
                   };

                   handleAddTab(newTab);
                   setIsAddingTab(false);
                   setNewTabText('');
                   e.currentTarget.value = '';
               }
        }
    }

    const handleTabPlace = (movedTabIndex: number) => {
        const movedTab = tabs[movedTabIndex];
        const movedTabElement = tabsRef.current[movedTab.tabRoute].current;
        if(!movedTabElement) return;
        const movedTabRect = movedTabElement.getBoundingClientRect()

        let furthestPassedIndex = 0;
        for(let i = 1; i < tabs.length; i++) {
            const tab = tabs[i];
            const tabElement = tabsRef.current[tab.tabRoute];
            if(!tabElement.current) continue;

            const tabRect = tabElement.current.getBoundingClientRect();
            
            if(movedTabRect.x < tabRect.x) break;

            furthestPassedIndex = i;
        }

        if(furthestPassedIndex === movedTabIndex) return;

        let leftIndex: number = 0;
        let rightIndex: number = 0;
        let dir: 'left'|'right' = 'left';
        if(movedTabIndex > furthestPassedIndex) {
            leftIndex = Math.min(furthestPassedIndex + 1, tabs.length - 1);
            rightIndex = movedTabIndex;
            dir = 'right';
        } else if(movedTabIndex < furthestPassedIndex) {
            leftIndex = movedTabIndex;
            rightIndex = furthestPassedIndex;
            dir = 'left';
        }

        rotateTabs(leftIndex, rightIndex, dir);
    }

    const listenersOn = () => {
        const offFns: (() => void)[] = [];
        offFns.push(Events.On('common:shared-chat-begin', handleSharedChatBegin));
        offFns.push(Events.On('common:shared-chat-update', handleSharedChatUpdate));
        offFns.push(Events.On('common:shared-chat-end', handleSharedChatEnd));
    }

    useEffect(() => {
        if(location.hash.slice(1) !== currentTabRoute) navigate(currentTabRoute);

        return listenersOn();
    }, []);

    return (
        <div className={'flex max-w-full items-center w-full h-9 gap-1 border-b p-1 border-outline-2 bg-bg-09'}
        >
            <div className={'shrink-0 flex items-center justify-center w-7 h-7 border border-outline-1 rounded-sm p-1 [&_svg]:fill-text-1 data-[selected=true]:bg-bg-5 hover:bg-bg-8 cursor-pointer'}
                onClick={() => handleTabSelect(home)}
                data-selected={currentTabRoute === home.tabRoute ? 'true' : 'false'}
            >
                <HomeIcon />
            </div>
            {tabs.slice(1).map((tab, i) =>
            <Tab
                tab={tab}
                index={i+1}
                key={tab.tabRoute}
                ref={tabsRef.current[tab.tabRoute] ??= { current: null }}
                selected={tab.tabRoute === currentTabRoute}
                onTabSelect={handleTabSelect}
                onTabRemove={handleTabRemove}
                onTabPlace={() => handleTabPlace(i+1)}
            />
                     )}
            <div className={(isAddingTab ? '' : ' hidden') + ' flex justify-start items-center border border-text-1 rounded-sm text-sm p-1 h-7 bg-bg-09 max-w-50 basis-25'}
            >
                <input
                    className="max-w-50 min-w-25 h-full bg-transparent! border-none!"
                    type='text'
                    ref={(node: HTMLInputElement|null) => {node && node.focus()}}
                    onKeyDown={handleAddTabKeyDown}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => setNewTabText(e.currentTarget.value)}
                    onBlur={() => setIsAddingTab(false)}
                />
            </div>
            <div
                className={'shrink-0 flex justify-center items-center w-7 h-7 p-2 [&_svg]:fill-text-1 hover:[&_svg]:fill-text-2 hover:[&_svg]:brightness-90 cursor-pointer' + (isAddingTab ? ' hidden' : '')}
                onClick={() => setIsAddingTab(true)}
            >
                <PlusIcon />
            </div>
            <div
                className={'shrink-0 flex justify-center items-center w-7 h-7 p-2 [&_svg]:fill-text-1 hover:[&_svg]:fill-text-2 hover:[&_svg]:brightness-90 cursor-pointer' + (isAddingTab ? ' hidden' : '')}
                onClick={() => navigate('/search')}
            >
                <SearchIcon />
            </div>
        </div>
    )
}
