export const U_NBSP = 160; // non-breaking space
export const U_SP = 32;    // space

export function moveCursorToEnd(element: HTMLElement) {

    if(element.nodeName !== 'textarea' &&
       element.getAttribute('contenteditable') !== 'true'
      ) return;

    const length = element.childNodes.length;
    window.getSelection()?.setPosition(element, length);
}

export function getCursorPos(element: HTMLElement): number {
    const selection = window.getSelection();
    if(!selection || selection.rangeCount === 0 || !element.contains(selection.anchorNode)) return -1;

    const range = selection.getRangeAt(0);
    const clRange = range.cloneRange();
    clRange.selectNodeContents(element)
    clRange.setEnd(range.endContainer, range.endOffset);
    const i = Array.from(element.childNodes.values()).indexOf(range.endContainer as ChildNode);

    return i + clRange.toString().length;
}

export function moveCursorTo(element: HTMLElement, pos: number) {
    element.focus();
    const selection = window.getSelection();
    if(!selection) return;

    let curPos = 0;
    for(let i = 0; i < element.childNodes.length; i++, curPos++) {
        const childNode = element.childNodes.item(i);
        if(childNode.nodeType === Node.TEXT_NODE && childNode.nodeValue) {
            curPos += childNode.nodeValue.length;
       }
       if(curPos >= pos) {
            const offset = curPos - pos;
            const range = document.createRange();

            range.setStart(childNode, offset);
            range.setEnd(childNode, offset);

            selection.removeAllRanges();
            selection.addRange(range);
            return;
       }
    }

    moveCursorToEnd(element);
}
