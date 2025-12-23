export class DebugLogger {
    funcLog(func: Function, message: string) {
        console.log(`[${func.name}]: ${message}`);
    }

    funcErr(func: Function, err: string) {
        console.error(`[${func.name}]: ${err}`);
    }
}
