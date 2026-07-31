import { isDefined } from '@/util/assert';
import React, { isValidElement, ReactElement, useEffect, useRef, useState } from 'react';

export interface IJumpToRecentPopupProps extends React.HTMLAttributes<HTMLDivElement>{
    newItems: number;
}

interface IAutoScrollerProps extends React.HTMLAttributes<HTMLDivElement> {
    jumpToRecentPopup?: React.FC<IJumpToRecentPopupProps>;
    scrollThresh?: number,
    children?: React.ReactNode;
}

export default function AutoScroller({
    jumpToRecentPopup,
    scrollThresh=200,
    children,
    ...rest
}: IAutoScrollerProps) {
    const [atBottom, setAtBottom] = useState<boolean>(true);
    const [newItems, setNewItems] = useState<number>(0);
    const anchorRef = useRef<HTMLDivElement|null>(null);
    const prevLastChild = useRef<ReactElement|null>(null);

    const scrollerRef = useRef<HTMLDivElement>(null);


    const scrollToBottom = (behavior: ScrollBehavior = "smooth") => {
        anchorRef.current?.scrollIntoView({behavior})
    }

    const shouldAutoScroll = (element: HTMLElement) => {
        const { scrollHeight, scrollTop, clientHeight } = element;
        const distanceFromBottom = scrollHeight - scrollTop - clientHeight;

        return distanceFromBottom <= scrollThresh;
    }

    const handleScroll = () => {
        if(!scrollerRef.current) return;

        setAtBottom(shouldAutoScroll(scrollerRef.current));
    }

    useEffect(() => {
        if(!Array.isArray(children)) return;

        if(atBottom) {
            scrollToBottom('instant');
            setNewItems(0);
            prevLastChild.current = children.at(-1);
        } else if(newItems <= 10 && isDefined(prevLastChild.current)) {
            const msgKey = prevLastChild.current.key;
            let offset = 0;
            for(; offset <= 10; offset++) {
                const child = children.at(children.length-offset-1);
                if(isValidElement(child) && msgKey === child.key) {
                    break;
                }
            }

            setNewItems(offset);
        }
    }, [children]);

    return (
        <>
            <div
            className='w-full h-full flex flex-col scroller-y z-200'
            ref={scrollerRef}
            onScroll={handleScroll}
            {...rest}
            >
                {children}
            <div className="h-0 w-0" ref={anchorRef}></div>
            </div>
            <div className={`w-full absolute flex justify-center items-center ${atBottom ? 'invisible' : 'visible'}`}
                style={{
                    top: `calc(${scrollerRef.current?.offsetHeight.toString().concat('px') ?? '100%'} * 0.925)`
                }}
            >
                {jumpToRecentPopup?.({
                    newItems,
                    onClick: () => scrollToBottom(),
                    style: {
                        zIndex: 300,
                    }
                }) as React.ReactNode}
            </div>
        </>
    )
}
