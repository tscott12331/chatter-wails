import styles from './jump-to-recent-popup.module.css';

import { IJumpToRecentPopupProps } from '../util/auto-scroller';

export default function JumpToRecentPopup({
    newItems,
    ...rest
}: IJumpToRecentPopupProps) {
    const MAX_ITEM_DISPLAY = 10;

    const generatePopupText = (numItems): string => {
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
        className={styles.wrapper + ' flex-center'}
        {...rest}
        >
            <p>{generatePopupText(newItems)}</p>
        </div>
    )
}
