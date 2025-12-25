import { TTab } from './tab-manager';
import PlusIcon from '../svg/plus-icon';
import { useEffect, useState } from 'react';

interface ITabProps {
    tab: TTab;
    index: number;
    selected: boolean;
    onTabSelect: (tab: TTab) => void;
    onTabRemove: (tab: TTab) => void;
    onTabMove: (tab: TTab, x: number) => void;
    ref: React.RefObject<HTMLDivElement|null>;
}

export default function Tab({
    tab,
    index,
    selected,
    onTabSelect,
    onTabRemove,
    onTabMove,
    ref,
}: ITabProps) {
    const [isDragging, setIsDragging] = useState<boolean>(false);
    const [mouseX, setMouseX] = useState<number>(0);
    const [initX, setInitX] = useState<number>(0);

    const handleTabRemove = (e: React.MouseEvent<HTMLDivElement>) => {
        e.stopPropagation();
        onTabRemove(tab);
    }

    const handleMouseDown = (e: React.MouseEvent<HTMLDivElement>) => {
        onTabSelect(tab);
        setInitX(e.pageX);
        setMouseX(e.pageX);
        setIsDragging(true);
    }

    const handleMouseUp = (e: MouseEvent) => {
        setIsDragging(false);
    }

    const handleMouseMove = (e: MouseEvent) => {
        setMouseX(e.pageX);
        onTabMove(tab, e.pageX);
    }

    useEffect(() => {
        if(isDragging) {
            document.addEventListener('mousemove', handleMouseMove);

            return () => document.removeEventListener('mousemove', handleMouseMove);
        }

        return;
    }, [isDragging]);

    useEffect(() => {
        document.addEventListener('mouseup', handleMouseUp);

        return () => document.removeEventListener('mouseup', handleMouseUp);
    }, []);

    useEffect(() => {
        if(ref.current) {
            setInitX(mouseX);
        }
    }, [index]);

    return (
        <div
        className={`${selected ? 'bg-bg-5' : ''} max-w-50 min-w-25 border border-text-1 rounded-sm text-sm p-1 h-7 gap-1 select-none flex items-center justify-between relative hover:bg-bg-8 hover:[&_svg]:inline-block `}
        onMouseDown={handleMouseDown}
        ref={ref}
        style={{
            transform: `translateX(${isDragging ? `${mouseX - initX}px` : '0'})`,
            zIndex: isDragging ? '5000' : 'initial',
        }}
        >
        <p className='text-ellipsis grow-2'>{tab.tabName}</p>
        <div
        className={'w-3 h-3 flex justify-center items-center [&_svg]:fill-red-200 hover:[&_svg]:fill-red-300 [&_svg]:rotate-45 [&_svg]:hidden'}
        onClick={handleTabRemove}
        >
            <PlusIcon />
        </div>
        </div>
    )
}
