export function moveCursorToEnd(element: HTMLElement) {

    if(element.nodeName !== 'textarea' &&
       element.getAttribute('contenteditable') !== 'true'
      ) return;

    const length = element.childNodes.length;
    window.getSelection()?.setPosition(element, length);
}
