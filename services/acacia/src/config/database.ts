import { PostgresDialect, Kysely } from 'kysely';
import { Pool } from 'pg';
import environment from './environment.ts';
import type { DB } from '#dashboard-api/types/database.d.ts';

const dialect = new PostgresDialect({
  pool: new Pool({ connectionString: environment.database_url, max: 10 }),
});

const database = new Kysely<DB>({ dialect });

export { database, dialect };
