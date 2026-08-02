import { TTab } from "@/components/tabs/tab-manager";
import { createContext, useContext, useMemo, useState } from "react";
import { GlobalContext } from "./global-context";
import { DisconnectFromChatroom } from "@wailsjs/chatter-wails/appservice";
import { rotateArr } from "@/util/arr";

interface ITabContext {
    home: TTab;
    tabs: TTab[];
    addTab: (tab: TTab) => void;
    removeTab: (tab: TTab) => void;
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
    addTab(_tab) {},
    removeTab(_tab) {},
    editTab(_index, _changeFn) {},
    rotateTabs(_leftIndex, _rightIndex, _dir) {},
});

export function TabContextProvider({
    children
}: { children: React.ReactNode }) {

    const { broadcastError } = useContext(GlobalContext);

    const [tabs, setTabs] = useState<TTab[]>([HOME_TAB]);

    const removeTab = (tab: TTab) => {
        if(tab.tabRoute.toLowerCase() == HOME_TAB.tabRoute.toLowerCase()) return; // can't remove home tab :)

        DisconnectFromChatroom(tab.tabRoute.split('/chatroom/')[1]).catch(broadcastError);

        setTabs((curTabs) => curTabs.filter(t => t.tabRoute !== tab.tabRoute));
    }

    const addTab = (tab: TTab) => {
        setTabs((curTabs) => [...curTabs, tab]);
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
        addTab,
        removeTab,
        editTab,
        rotateTabs,
    }), [tabs])

    return (
        <TabContext.Provider value={ctxValue}>
            {children}
        </TabContext.Provider>
    )

}
