import React, { useEffect, useRef, useState } from 'react';

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
    const [shouldShowPopup, setShouldShowPopup] = useState<boolean>(false);
    const [newItems, setNewItems] = useState<number>(0);

    /*
    Children gets updated in useEffect every time
    popup is opened.
    Keeping track of the previous last child prevents
    us incrementing the newItems count when the last item
    doesn't change
    */
    const [prevChild, setPrevChild] = useState<Element|null>(null);

    const scrollerRef = useRef<HTMLDivElement>(null);


    const scrollToBottom = (behavior: ScrollBehavior = "smooth") => {
        scrollerRef.current?.scroll({
            behavior,
            top: scrollerRef.current.scrollHeight
        })
    }

    const handleItemScroll = (e: React.UIEvent<HTMLDivElement>) => {
        if(scrollerRef.current) {
            const scrollHeight = e.currentTarget.scrollHeight;
            const scrollTop = e.currentTarget.scrollTop;
            const wrapperHeight = e.currentTarget.offsetHeight;

            const passedThresh = scrollHeight - (scrollTop + wrapperHeight) >= scrollThresh;
            setShouldShowPopup(passedThresh);
            if(!shouldShowPopup) {
                setNewItems(0);
            }
        }
    }

    useEffect(() => {
        if(!scrollerRef.current) return;
        if(Array.isArray(children) && children.length > 0) {
            const scrollHeight = scrollerRef.current.scrollHeight;
            const scrollTop = scrollerRef.current.scrollTop;
            const wrapperHeight = scrollerRef.current.offsetHeight;
            const lastChild = scrollerRef.current.children.item(children.length-1);
            if(!lastChild || lastChild === prevChild) return;
            setPrevChild(lastChild);

            const itemHeight = lastChild.clientHeight;

            if((scrollHeight - itemHeight) - (scrollTop + wrapperHeight) < scrollThresh) {
                  setShouldShowPopup(false);
                  scrollToBottom('instant');
                  setNewItems(0);
              } else {
                  setShouldShowPopup(true);
                  setNewItems((items) => items + 1);
              }
        }

    }, [children]);
    return (
        <>
            <div
            className='w-full h-full flex flex-col scroller-y z-200'
            ref={scrollerRef}
            onScroll={handleItemScroll}
            {...rest}
            >
                {children}
            </div>
            <div className={`w-full absolute flex justify-center items-center ${shouldShowPopup ? 'visible' : 'invisible'}`}
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
