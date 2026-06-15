export function listToMap<T>(list: T[], key: keyof T): Map<T[keyof T], T> {
    return new Map(list.map(e => [e[key], e]));
}
