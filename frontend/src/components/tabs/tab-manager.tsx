import { useEffect, useRef, useState } from 'react';
import Tab from './tab';
import HomeIcon from '../svg/home-icon';
import PlusIcon from '../svg/plus-icon';
import { useNavigate } from 'react-router-dom';
import { rotateArr } from '@util/arr';

export type TTab = {
    tabRoute: string;
    tabName: string;
}

interface ITabManagerProps {

}

export default function TabManager({

}: ITabManagerProps) {
    const HOME_TAB: TTab = {
        tabRoute: '/',
        tabName: 'home',
    };

    const [tabs, setTabs] = useState<TTab[]>([HOME_TAB]);
    const [currentTabRoute, setCurrentTabRoute] = useState<string>('/');
    const tabsRef = useRef<Record<string, {current: HTMLDivElement|null}>>({})

    const [isAddingTab, setIsAddingTab] = useState<boolean>(false);
    const [newTabText, setNewTabText] = useState<string>("");

    const navigate = useNavigate();

    const handleTabSelect = (tab: TTab) => {
        setCurrentTabRoute(tab.tabRoute);

        navigate(tab.tabRoute);
    }

    const handleTabRemove = (tab: TTab) => {
        if(tab.tabRoute == HOME_TAB.tabRoute) return; // can't remove home tab :)

        if(tab.tabRoute === currentTabRoute) {
            setCurrentTabRoute(HOME_TAB.tabRoute);
            navigate(HOME_TAB.tabRoute);
        }

        setTabs((curTabs) => curTabs.filter(t => t.tabRoute !== tab.tabRoute));
        delete tabsRef.current[tab.tabRoute];
    }

    const handleAddTab = (tab: TTab) => {
        if(tabs.find(t => t.tabRoute === tab.tabRoute)) {
            return handleTabSelect(tab)
        }

        setTabs((curTabs) => [...curTabs, tab]);
    }

    const handleAddTabKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if(e.key === 'Enter') {
            const trimmedTabName = newTabText.trim();
            if(!trimmedTabName.includes(' ') && trimmedTabName.length >= 3
               && trimmedTabName.length <= 25) {
                   const newTab: TTab = {
                        tabRoute: `/chatroom/${trimmedTabName.toLowerCase()}`,
                        tabName: trimmedTabName
                   };
                   handleAddTab(newTab);
                   setIsAddingTab(false);
                   setNewTabText('');
                   e.currentTarget.value = '';
               }
        }
    }

    const handleTabMove = (movedTab: TTab, movedTabX: number) => {
        const movedTabIndex = tabs.findIndex(t => t.tabRoute === movedTab.tabRoute);
        if(movedTabIndex === -1) return;

        const movedTabElement = tabsRef.current[movedTab.tabRoute].current;
        if(!movedTabElement) return;
        const movedTabRect = movedTabElement.getBoundingClientRect()
        // console.log(`rect.x: ${movedTabRect.x}`);
        // console.log(`movedTabDX: ${movedTabDX}`);
        // console.log(`movedTabX: ${movedTabX}`);

        const entries = Object.entries(tabsRef.current)
                            .filter(([t, _]) => t != movedTab.tabRoute);

        

        let furthestPassedTabRoute = HOME_TAB.tabRoute;
        let logString = "\nSTART\n";
        for(const [i, [tabRoute, element]] of entries.entries()) {
            const rect = element.current?.getBoundingClientRect()
            logString += `tab index ${i} ${tabRoute}\n`
            if(!rect) continue;
            const middleOfTab = rect.x + rect.width/2;
            const middleOfMovedTab = movedTabRect.x + movedTabRect.width/2;
            logString += `middleOfTab: ${middleOfTab}\n`;
            logString += `middleOfMovedTab: ${middleOfMovedTab}\n`;

            // tabs are stored from left to right
            // keep updating furthestPassedTab until tabX > x
            if(middleOfMovedTab < middleOfTab) {
                break;
            }

            furthestPassedTabRoute = tabRoute;
        }
        // console.log(logString);

        const furthestPassedTabIndex = tabs.findIndex(t => t.tabRoute === furthestPassedTabRoute);
        if(furthestPassedTabIndex === -1) return;

        let leftIndex: number = 0;
        let rightIndex: number = 0;
        let dir: 'left'|'right' = 'left';
        if(movedTabIndex > furthestPassedTabIndex) {
            leftIndex = Math.min(furthestPassedTabIndex + 1, tabs.length - 1);
            rightIndex = movedTabIndex;
            dir = 'right';
        } else if(movedTabIndex < furthestPassedTabIndex) {
            leftIndex = movedTabIndex;
            rightIndex = furthestPassedTabIndex;
            dir = 'left';
        } else {
            return;
        }

        if(leftIndex == rightIndex) return;

        const preRotateSegement = [...tabs.slice(leftIndex, rightIndex + 1)];
        const rotatedSegement = rotateArr(preRotateSegement, dir)
        const newTabs = [...tabs.slice(0, leftIndex),
                 ...rotatedSegement,
                 ...tabs.slice(rightIndex + 1)
                ];
        setTabs(newTabs);
    }

    useEffect(() => {
        if(location.hash.slice(1) !== currentTabRoute) navigate(currentTabRoute);
    }, []);

    return (
        <div className={'flex items-center w-full h-9 gap-1 border-b p-1 border-outline-2 bg-bg-09'}
        >
            <div className={'flex items-center justify-center w-7 h-7 border border-outline-1 rounded-sm p-1 [&_svg]:fill-text-1 data-[selected=true]:bg-bg-5 hover:bg-bg-8 cursor-pointer'}
                onClick={() => handleTabSelect(HOME_TAB)}
                data-selected={currentTabRoute === HOME_TAB.tabRoute ? 'true' : 'false'}
            >
                <HomeIcon />
            </div>
            {tabs.slice(1).map((tab, i) =>
            <Tab
                tab={tab}
                index={i}
                key={tab.tabRoute}
                ref={tabsRef.current[tab.tabRoute] ??= { current: null }}
                selected={tab.tabRoute === currentTabRoute}
                onTabSelect={handleTabSelect}
                onTabRemove={handleTabRemove}
                onTabMove={handleTabMove}
            />
                     )}
            <div className={(isAddingTab ? '' : ' hidden') + ' flex justify-start items-center border border-text-1 rounded-sm text-sm p-1 h-7 bg-bg-09 max-w-50 min-w-25'}
            >
                <input
                    className="max-w-50 min-w-25 h-full bg-transparent! border-none!"
                    type='text'
                    ref={(node) => node && node.focus()}
                    onKeyDown={handleAddTabKeyDown}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => setNewTabText(e.currentTarget.value)}
                    onBlur={() => setIsAddingTab(false)}
                />
            </div>
            <div
                className={'flex justify-center items-center w-7 h-7 p-2 [&_svg]:fill-text-1 hover:[&_svg]:fill-text-2 hover:[&_svg]:brightness-90 cursor-pointer' + (isAddingTab ? ' hidden' : '')}
                onClick={() => setIsAddingTab(true)}
            >
                <PlusIcon />
            </div>
        </div>
    )
}
