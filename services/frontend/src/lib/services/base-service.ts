import { BaseLogger } from '../config/logger';
import { cookieService } from './cookie-service';

export type BaseServiceArguments = {
  url: string;
  logger: BaseLogger;
};

export abstract class BaseHttpService {
  protected url: string;
  protected logger: BaseLogger;
  protected cookieService: typeof cookieService;

  constructor({ url, logger }: BaseServiceArguments) {
    this.url = url;
    this.logger = logger;
    this.cookieService = cookieService;
  }
}
