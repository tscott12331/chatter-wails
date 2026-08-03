import { createContext, useContext, useMemo, useState } from "react";
import { GlobalContext } from "./global-context";
import { DisconnectFromChatroom } from "@wailsjs/chatter-wails/appservice";
import { rotateArr } from "@/util/arr";

export type TTab = {
    readonly tabRoute: string;
    readonly tabName: string;
}

export const createTabRoute = (channel: string): string => {
    return `/chatroom/${channel.toLowerCase()}`;
}

export const createTab = (channel: string): TTab => {
    return {
        tabName: channel,
        tabRoute: createTabRoute(channel),
    };
}

interface ITabContext {
    home: TTab;
    tabs: TTab[];
    curTab: TTab;
    addTab: (tab: TTab) => void;
    removeTab: (tab: TTab) => void;
    selectTab: (tab: TTab) => void;
    editTab: (index: number, changeFn: (tab: TTab) => TTab) => void;
    rotateTabs: (leftIndex: number, rightIndex: number, dir: 'left'|'right') => void;
}

const HOME_TAB: TTab = {
    tabRoute: '/',
    tabName: 'home',
};

export const TabContext = createContext<ITabContext>({
    home: HOME_TAB,
    tabs: [],
    curTab: HOME_TAB,
    addTab(_tab) {},
    removeTab(_tab) {},
    selectTab(_tab) {},
    editTab(_index, _changeFn) {},
    rotateTabs(_leftIndex, _rightIndex, _dir) {},
});

export function TabContextProvider({
    children
}: { children: React.ReactNode }) {

    const { broadcastError } = useContext(GlobalContext);

    const [tabs, setTabs] = useState<TTab[]>([HOME_TAB]);
    const [curTab, setCurTab] = useState<TTab>(HOME_TAB);

    const removeTab = (tab: TTab) => {
        if(tab.tabRoute == HOME_TAB.tabRoute) return; // can't remove home tab :)
        if(tab.tabRoute === curTab.tabRoute) {
            setCurTab(HOME_TAB);
        }

        DisconnectFromChatroom(tab.tabRoute.split('/chatroom/')[1]).catch(broadcastError);

        setTabs((curTabs) => curTabs.filter(t => t.tabRoute !== tab.tabRoute));
    }

    const addTab = (tab: TTab) => {
        if(!tabs.find(t => t.tabRoute === tab.tabRoute)) {
            setTabs((curTabs) => [...curTabs, tab]);
        }
    }

    const selectTab = (tab: TTab) => {
        if(tabs.find(t => t.tabRoute === tab.tabRoute)) {
            setCurTab(tab);
        }
    }

    const editTab = (index: number, changeFn: (tab: TTab) => TTab) => {
        setTabs(tabs => {
            if(index < 0 || index >= tabs.length) return tabs;
            const changed = changeFn(tabs[index]);
            const newTabs = [
                ...tabs.slice(0, index),
                changed,
                ...tabs.slice(index+1)
            ];
            return newTabs;
        });
    }

    const rotateTabs = (leftIndex: number, rightIndex: number, dir: 'left'|'right') => {
        if(leftIndex == rightIndex) return;

        const preRotateSegement = [...tabs.slice(leftIndex, rightIndex + 1)];
        const rotatedSegement = rotateArr(preRotateSegement, dir)
        const newTabs = [
            ...tabs.slice(0, leftIndex),
            ...rotatedSegement,
            ...tabs.slice(rightIndex + 1)
        ];
        setTabs(newTabs);
    }

    const ctxValue = useMemo<ITabContext>(() => ({
        home: HOME_TAB,
        tabs,
        curTab,
        addTab,
        removeTab,
        selectTab,
        editTab,
        rotateTabs,
    }), [tabs, curTab])

    return (
        <TabContext.Provider value={ctxValue}>
            {children}
        </TabContext.Provider>
    )

}
