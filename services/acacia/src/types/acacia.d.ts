import type { AcaciaDatabaseType } from '#acacia/config/database.ts';

declare global {
  namespace Express {
    interface Request {
      acacia: {
        db: AcaciaDatabaseType;
      };
    }
  }
}
