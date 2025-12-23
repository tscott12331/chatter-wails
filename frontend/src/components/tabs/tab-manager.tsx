import styles from './tab-manager.module.css';

import { useEffect, useRef, useState } from 'react';
import Tab from './tab';
import HomeIcon from '../svg/home-icon';
import PlusIcon from '../svg/plus-icon';
import { useNavigate } from 'react-router-dom';
import { deepEqual } from '../../util/obj'
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

    const addingTabRef = useRef<HTMLInputElement>(null);
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
                   if(addingTabRef.current) addingTabRef.current.value = '';
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

    useEffect(() => {
        addingTabRef.current?.focus();
    }, [isAddingTab]);

    useEffect(() => {
        navigate(currentTab);
    }, []);

    return (
        <div className={styles.wrapper + ' flex-align-center'}
        >
            <div className={styles.homeTab + ' flex-center'}
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
            <div className={styles.addTabPopup + (isAddingTab ? '' : ' no-display') + ' flex-justify-start flex-align-center'}
            >
                <input
                    type='text'
                    ref={addingTabRef}
                    onKeyDown={handleAddTabKeyDown}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => setNewTabText(e.currentTarget.value)}
                    onBlur={() => setIsAddingTab(false)}
                />
            </div>
            <div
                className={styles.addTab + ' flex-center' + (isAddingTab ? ' no-display' : '')}
                onClick={() => setIsAddingTab(true)}
            >
                <PlusIcon />
            </div>
        </div>
    )
}
