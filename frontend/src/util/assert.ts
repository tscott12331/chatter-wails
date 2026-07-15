export function assert(condition: any, message?: string): asserts condition {
    if(!condition) {
        throw new Error(message || "Assertion error");
    }
}

export function assertDefined<T>(value: T, message?: string): asserts value is NonNullable<T> {
    assert(isDefined(value), message);
}

export function isDefined<T>(value: T): value is NonNullable<T> {
    return value !== null && value !== undefined;
}
