import { IJumpToRecentPopupProps } from '../util/auto-scroller';

export default function JumpToRecentPopup({
    newItems,
    ...rest
}: IJumpToRecentPopupProps) {
    const MAX_ITEM_DISPLAY = 10;

    const generatePopupText = (numItems: number): string => {
        if(numItems > MAX_ITEM_DISPLAY) {
            return String(MAX_ITEM_DISPLAY).concat('+ new messages');
        } else if(numItems > 0) {
            return `${numItems} new message${numItems > 1 ? 's' : ''}`
        } else {
            return 'Jump to new messages';
        }
    }

    return (
        <div
        className='flex justify-center items-center p-1.5 bg-bg-3/50 backdrop-blur-xs rounded-xs cursor-pointer'
        {...rest}
        >
            <p>{generatePopupText(newItems)}</p>
        </div>
    )
}
