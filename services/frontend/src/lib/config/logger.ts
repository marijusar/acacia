export interface BaseLogger {
  info: (message: string, context?: object) => void;
  warn: (message: string, context?: object) => void;
  error: (message: string, context?: object) => void;
  debug: (message: string, context?: object) => void;
}

class Logger implements BaseLogger {
  info(message: string, context?: object) {
    console.info(message, context);
  }

  warn(message: string, context?: object) {
    console.warn(message, context);
  }

  error(message: string, context?: object) {
    console.error(message, context);
  }

  debug(message: string, context?: object) {
    console.debug(message, context);
  }
}

export const logger = new Logger();