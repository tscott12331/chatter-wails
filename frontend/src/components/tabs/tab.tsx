import { TTab } from '@/contexts/tab-context';
import PlusIcon from '../svg/plus-icon';
import { useEffect, useState } from 'react';

interface ITabProps {
    tab: TTab;
    index: number;
    selected: boolean;
    onTabSelect: (tab: TTab) => void;
    onTabRemove: (tab: TTab) => void;
    onTabPlace: (tab: TTab) => void;
    ref?: React.RefObject<HTMLDivElement|null>;
}

export default function Tab({
    tab,
    selected,
    onTabSelect,
    onTabRemove,
    onTabPlace,
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

    return (
        <div
        className={`${selected ? 'bg-bg-5' : ''} ${isDragging ? 'bg-bg-5/70 backdrop-blur-xs' : 'hover:bg-bg-6'} shrink grow basis-25 min-w-0 max-w-50 border border-outline-1 rounded-sm text-sm p-1 h-7 gap-1 select-none flex items-center justify-between relative hover:[&_svg]:inline-block cursor-pointer`}
        onMouseDown={handleMouseDown}
        onMouseUp={() => onTabPlace(tab)}
        ref={ref}
        style={{
            transform: `translateX(${isDragging ? `${mouseX - initX}px` : '0'})`,
            zIndex: isDragging ? '5000' : 'initial',
        }}
        >
            <p className='text-ellipsis whitespace-nowrap overflow-hidden grow-2'>{tab.tabName}</p>
            <div
            className={'w-3 h-3 flex justify-center items-center [&_svg]:fill-red-200 hover:[&_svg]:fill-red-300 [&_svg]:rotate-45 [&_svg]:hidden'}
            onClick={handleTabRemove}
            onMouseDown={(e) => e.stopPropagation()}
            onMouseUp={(e) => e.stopPropagation()}
            >
                <PlusIcon />
            </div>
        </div>
    )
}
