import styles from './tooltip.module.css';
import React, { useState } from "react";

interface ITooltipProps extends React.HTMLAttributes<HTMLDivElement> {
    text: string;
    hoverTime?: number;
    children: React.ReactNode;
}

export default function Tooltip({
 text,
 hoverTime=1000,
 children,
 ...rest
}: ITooltipProps) {
    const POPUP_POS_OFFSET = 15;

    const [shouldShow, setShouldShow] = useState<boolean>(false);
    const [pos, setPos] = useState<{x: number, y: number}>({x: 0, y: 0});

    let timeout: NodeJS.Timeout;

    const showPopup = (x: number, y: number) => {
        setShouldShow(true);
        setPos({
            x,
            y,
        });
    }

    const handleMouseEnter = (e: React.MouseEvent<HTMLDivElement>) => {
        const x = e.pageX / e.currentTarget.clientWidth + POPUP_POS_OFFSET;
        const y = e.pageY / e.currentTarget.clientHeight + POPUP_POS_OFFSET;
        if(hoverTime > 0) {
            handleMouseLeave();
            timeout = setTimeout(() => {
                showPopup(x, y);
            }, hoverTime);
        } else {
            showPopup(x, y);
        }

    }

    const handleMouseLeave = () => {
        clearTimeout(timeout);

        setShouldShow(false);
    }

    return (
        <span
            onMouseEnter={handleMouseEnter}
            onMouseLeave={handleMouseLeave}
            style={{position: 'relative'}}
            {...rest}
        >
            {children}
            {shouldShow &&
            <span
                className={styles.popupWrapper + ' absolute z-2000'}
                style={{
                    left: `${pos.x}px`,
                    top: `${pos.y}px`,
                }}
            >{text}</span>
            }
        </span>
    )
}
