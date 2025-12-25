import { useRef, useState } from 'react';
import Tab from './tab';
import HomeIcon from '../svg/home-icon';
import PlusIcon from '../svg/plus-icon';
import { useNavigate } from 'react-router-dom';
import { deepEqual } from '@util/obj'
import { rotateArr } from '@util/arr';

export type TTab = {
    tab: string;
    tabName: string;
}

interface ITabManagerProps {

}

export default function TabManager({

}: ITabManagerProps) {
    const HOME_TAB: TTab = {
        tab: '/',
        tabName: 'home',
    };

    const [tabs, setTabs] = useState<TTab[]>([HOME_TAB]);
    const [currentTab, setCurrentTab] = useState<string>('/');
    const tabsRef = useRef<Record<string, {current: HTMLDivElement|null}>>({})

    const [isAddingTab, setIsAddingTab] = useState<boolean>(false);
    const [newTabText, setNewTabText] = useState<string>("");

    const navigate = useNavigate();

    const handleTabSelect = (tab: TTab) => {
        setCurrentTab(tab.tab);

        navigate(tab.tab);
    }

    const handleTabRemove = (tab: TTab) => {
        if(deepEqual(tab, HOME_TAB)) return; // can't remove home tab :)

        if(tab.tab === currentTab) {
            setCurrentTab(HOME_TAB.tab);
            navigate(HOME_TAB.tab);
        }

        setTabs((curTabs) => curTabs.filter(t => t.tab !== tab.tab));
    }

    const handleAddTab = (tab: TTab) => {
        if(tabs.includes(tab)) return handleTabSelect(tab);

        setTabs((curTabs) => [...curTabs, tab]);
    }

    const handleAddTabKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if(e.key === 'Enter') {
            const trimmedTabName = newTabText.trim();
            if(!trimmedTabName.includes(' ') && trimmedTabName.length >= 3
               && trimmedTabName.length <= 25) {
                   const newTab: TTab = {
                        tab: `/chatroom/${trimmedTabName}`,
                        tabName: trimmedTabName
                   };
                   handleAddTab(newTab);
                   setIsAddingTab(false);
                   setNewTabText('');
                   e.currentTarget.value = '';
               }
        }
    }

    const handleTabMove = (tab: TTab, x: number) => {
        const movingIndex = tabs.findIndex(t => t.tab === tab.tab);
        let indexOffset = 1;
        let comparedTab: HTMLDivElement|null|undefined;

        if(movingIndex - indexOffset > 0 &&
          (comparedTab = tabsRef.current[tabs[movingIndex - indexOffset].tab].current) &&
            x <= comparedTab?.offsetLeft + comparedTab.offsetWidth) {

            indexOffset++;
            while(movingIndex - indexOffset > 0 &&
              (comparedTab = tabsRef.current[tabs[movingIndex - indexOffset].tab].current) &&
                x <= comparedTab?.offsetLeft + comparedTab.offsetWidth) {
                indexOffset++;
            }

            setTabs([...tabs.slice(0, movingIndex - indexOffset + 1),
                    ...rotateArr(tabs.slice(movingIndex - indexOffset + 1, movingIndex + 1), 'right', 1),
                    ...tabs.slice(movingIndex + 1)]);
        } else if(movingIndex + indexOffset < tabs.length &&
                  (comparedTab = tabsRef.current[tabs[movingIndex + indexOffset].tab].current) &&
                     x >= comparedTab.offsetLeft) {

            indexOffset++;
            while(movingIndex + indexOffset < tabs.length &&
                  (comparedTab = tabsRef.current[tabs[movingIndex + indexOffset].tab].current) &&
                     x >= comparedTab.offsetLeft) {
                indexOffset++;
            }

            setTabs([...tabs.slice(0, movingIndex),
                    ...rotateArr(tabs.slice(movingIndex, movingIndex + indexOffset), 'right', indexOffset - 1),
                    ...tabs.slice(movingIndex + indexOffset)]);
        }
    }

    if(location.hash.slice(1) !== currentTab) navigate(currentTab);

    return (
        <div className={'flex items-center w-full h-9 gap-1 border-b p-1 border-outline-2'}
        >
            <div className={'flex items-center justify-center w-7 h-7 border border-text-1 rounded-sm p-1 [&_svg]:fill-text-1 data-[selected=true]:bg-bg-5 hover:bg-bg-8'}
                onClick={() => handleTabSelect(HOME_TAB)}
                data-selected={currentTab === HOME_TAB.tab ? 'true' : 'false'}
            >
                <HomeIcon />
            </div>
            {tabs.slice(1).map((tab, i) =>
            <Tab
                tab={tab}
                index={i}
                key={tab.tab}
                ref={tabsRef.current[tab.tab] ??= { current: null }}
                selected={tab.tab === currentTab}
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
                className={'flex justify-center items-center w-7 h-7 p-2 [&_svg]:fill-text-1 [&_svg]:hover:fill-text-2 [&_svg]:hover:brightness-90 [&_svg]:hover:drop-shadow-sm [&_svg]:hover:drop-shadow-outline-1 ' + (isAddingTab ? ' hidden' : '')}
                onClick={() => setIsAddingTab(true)}
            >
                <PlusIcon />
            </div>
        </div>
    )
}
