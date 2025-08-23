import { PostgresDialect, Kysely } from 'kysely';
import { Pool } from 'pg';
import type { DB } from '#acacia/types/database.d.ts';

const createDatabase = (url: string) => {
  const dialect = new PostgresDialect({
    pool: new Pool({ connectionString: url, max: 10 }),
  });

  const database = new Kysely<DB>({ dialect });

  return database;
};

export type AcaciaDatabaseType = ReturnType<typeof createDatabase>;

export { createDatabase };
