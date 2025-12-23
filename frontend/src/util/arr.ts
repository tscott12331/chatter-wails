export const rotateArr = (arr: any[], dir: 'left'|'right', amount: number = 1): any[] => {
    const newArr = [...arr];

    const addItem = dir === 'left' ? (item:any) => newArr.push(item) : (item: any) => newArr.unshift(item);
    const removeItem = dir === 'left' ? () => newArr.shift() : () => newArr.pop();

    for(let i = 0; i < amount; i++) {
        const item = removeItem();
        addItem(item);
    }

    return newArr;
}
