export const deepEqual = (obj1: object, obj2: object): boolean => {
    const obj1Keys = Object.keys(obj1);
    const obj2Keys = Object.keys(obj2);

    if(obj1Keys.length !== obj2Keys.length) return false

    for(const obj1Key of obj1Keys) {
        if(obj2Keys.includes(obj1Key)) {
            const obj1Value = obj1[obj1Key as keyof typeof obj1];
            const obj2Value = obj2[obj1Key as keyof typeof obj2];
            if(typeof obj1Value !== typeof obj2Value) return false;

            const valueType = typeof obj1Value;
            if(valueType === 'object') {
                if(!deepEqual(obj1Value, obj2Value)) return false;
            } else {
                if(obj1Value !== obj2Value) return false;
            }
        } else {
            return false;
        }
    }

    return true;
}
