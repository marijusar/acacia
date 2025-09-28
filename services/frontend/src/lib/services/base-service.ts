import { BaseLogger } from '../config/logger';

export type BaseServiceArguments = {
  url: string;
  logger: BaseLogger;
};

export class BaseHttpService {
  protected url: string;
  protected logger: BaseLogger;
  constructor({ url, logger }: BaseServiceArguments) {
    this.url = url;
    this.logger = logger;
  }
}
