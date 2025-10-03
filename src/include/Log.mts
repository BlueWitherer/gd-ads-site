/**
 * Colors and formatting for console logs
 */
export class LogFormat {
    public static gray = '\x1b[90m'; // debug
    public static white = '\x1b[37m'; // log
    public static cyan = '\x1b[36m'; // info
    public static yellow = '\x1b[93m'; // warn
    public static red = '\x1b[91m'; // error
    public static green = '\x1b[92m'; // done
    public static bold = '\x1b[1m'; // tag
    public static reset = '\x1b[0m'; // default
};

/**
 * Format time
 */
const timeStamp = (): string => {
    const now = new Date();

    const yyyy = now.getUTCFullYear();
    const mm = String(now.getUTCMonth() + 1).padStart(2, '0');
    const dd = String(now.getUTCDate()).padStart(2, '0');
    const h = String(now.getUTCHours()).padStart(2, '0');
    const m = String(now.getUTCMinutes()).padStart(2, '0');
    const s = String(now.getUTCSeconds()).padStart(2, '0');

    return `${yyyy}-${mm}-${dd} ${h}:${m}:${s} UTC`;
};

/**
 * Format multi-line logs
 * 
 * @param t Text
 * @param c Color code
 */
const formatLog = (t: string, c: LogFormat): string => {
    const lines = t.split('\n');
    return lines.map((ln) => `${c}${ln}${LogFormat.reset}`).join('\n');
};

/**
 * Get fully formatted log message
 * 
 * @param time Timestamp
 * @param color Color code
 * @param tag Log level
 * @param args All arguments
 */
const logMsg = (time: string, color: LogFormat, tag: string, ...args: [message?: any, ...optionalParams: any[]]): string => {
    const txt = args.join(' ');
    const msg = formatLog(txt, color);

    return `${time}${color} | ${LogFormat.bold}${tag}${LogFormat.reset}${color} | ${msg}${LogFormat.reset}`;
};

const safeParseLog = (err: unknown): string => {
    if (err) {
        try {
            if (typeof err === "object") {
                return JSON.stringify(err, null, 2);
            } else {
                return String(err);
            };
        } catch {
            return "[Unserializable log]";
        };
    } else {
        return "[Null log]";
    };
};

const formatArgs = (...args: any[]): string[] => {
    return args.map(safeParseLog);
};

const c = {
    debug: console.debug,
    info: console.info,
    warn: console.warn,
    error: console.error,
    log: console.log,
};

/**
 * Custom console methods with formatting
 */
export default class log {
    /**
     * Debug log
     * @param args 
     */
    public static debug = (...args: any[]): void => {
        c.debug(logMsg(timeStamp(), LogFormat.gray, 'DEBUG', formatArgs(...args)));
    };

    /**
     * Info log
     * @param args
     */
    public static info = (...args: any[]): void => {
        c.info(logMsg(timeStamp(), LogFormat.cyan, 'INFO', formatArgs(...args)));
    };

    /**
     * Warn log
     * @param args
     */
    public static warn = (...args: any[]): void => {
        c.warn(logMsg(timeStamp(), LogFormat.yellow, 'WARN', formatArgs(...args)));
    };

    /**
     * Error log
     * @param args
     */
    public static error = (...args: any[]): void => {
        c.error(logMsg(timeStamp(), LogFormat.red, 'ERROR', formatArgs(...args)));
    };

    /**
     * Done log
     * @param args
     */
    public static done = (...args: any[]): void => {
        c.log(logMsg(timeStamp(), LogFormat.green, 'DONE', formatArgs(...args)));
    };

    /**
     * Print log
     * @param args 
     */
    public static print = (...args: any[]): void => {
        c.log(logMsg(timeStamp(), LogFormat.white, ' LOG ', formatArgs(...args)));
    };
};

console.debug = (...args: [message?: any, ...optionalParams: any[]]) => {
    log.debug(...args);
};

console.info = (...args: [message?: any, ...optionalParams: any[]]) => {
    log.info(...args);
};

console.warn = (...args: [message?: any, ...optionalParams: any[]]) => {
    log.warn(...args);
};

console.error = (...args: [message?: any, ...optionalParams: any[]]) => {
    log.error(...args);
};

console.log = (...args: [message?: any, ...optionalParams: any[]]) => {
    log.print(...args);
};