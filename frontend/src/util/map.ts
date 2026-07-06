export function listToMap<T, K extends keyof T>(list: T[], key: K): Map<T[K], T> {
    return new Map(list.map(e => [e[key], e]));
}
