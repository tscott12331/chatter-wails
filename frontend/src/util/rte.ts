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
    const anchorNode = selection.anchorNode;
    if(!anchorNode) return -1;
    const range = selection.getRangeAt(0);

    let curPos = 0;
    for(let i = 0; i < element.childNodes.length; i++) {
        const childNode = element.childNodes.item(i);
        const isTextNode = childNode.nodeType === Node.TEXT_NODE;

        if(anchorNode?.isEqualNode(childNode)) {
            return curPos + range.startOffset;
        }

        let len = 0;
        if(isTextNode) {
            len = childNode.nodeValue?.length ?? 1;
        } else {
            len = 1;
        }

        curPos += len;
    }

    return -1;

    // const range = selection.getRangeAt(0);
    // const clRange = range.cloneRange();
    // clRange.selectNodeContents(element)
    // clRange.setEnd(range.endContainer, range.endOffset);
    // const i = Array.from(element.childNodes.values()).indexOf(range.endContainer as ChildNode);
    //
    // return i + clRange.toString().length;
}

export function moveCursorTo(element: HTMLElement, pos: number) {
    // console.log(`requested pos ${pos}`);
    element.focus();
    const selection = window.getSelection();
    if(!selection) return;

    let curPos = 0;
    for(let i = 0; i < element.childNodes.length; i++) {
        const childNode = element.childNodes.item(i);
        const isTextNode = childNode.nodeType === Node.TEXT_NODE;
        let len = 0;
        if(isTextNode) {
            len = childNode.nodeValue?.length ?? 1;
        } else {
            len = 1;
        }

        curPos += len;

        // console.log(`curPos: ${curPos}, len: ${len}`);

        if(curPos >= pos) {
            // console.log(`curPos is gte than seeked position`);
            if(!isTextNode) {
                // console.log(`not text, node, setting direct position to ${i}`)
                selection.setPosition(element, i + 1);
                return;
            }

            const diff = curPos - pos;
            const range = document.createRange();
            const offset = len - diff;
            // console.log(`offset: ${offset}`);

            // console.log(`requested pos within text node, offset ${offset}`);
            range.setStart(childNode, offset);
            range.setEnd(childNode, offset);

            selection.removeAllRanges();
            selection.addRange(range);
            return;
        }
    }

    moveCursorToEnd(element);
}
