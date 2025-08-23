import { PostgreSqlContainer } from '@testcontainers/postgresql';
import {
  FileMigrationProvider,
  Kysely,
  Migrator,
  PostgresDialect,
  sql,
} from 'kysely';
import { Pool } from 'pg';
import type { DB } from '#acacia/types/database.d.ts';
import { promises as fs } from 'node:fs';
import path from 'node:path';
import { randomUUID } from 'node:crypto';

export async function setupDatabaseContainer() {
  const container = await new PostgreSqlContainer('postgres:17')
    .withDatabase('acacia')
    .withUsername('postgres')
    .withPassword('root')
    .start();

  const connectionString = container.getConnectionUri();

  const dialect = new PostgresDialect({
    pool: new Pool({
      connectionString,
      max: 10,
    }),
  });

  const connection = new Kysely<DB>({ dialect });
  const migrator = new Migrator({
    db: connection,
    provider: new FileMigrationProvider({
      fs,
      path,
      // This needs to be an absolute path.
      migrationFolder: process.cwd() + '/migrations/',
    }),
  });

  // Wait for database to be ready
  await new Promise((resolve) => setTimeout(resolve, 2000));

  await migrator.migrateToLatest();

  return {
    container,
    connection,
    verify: () => !!container,
    createNewDatabase: async () => {
      const name = 'test_' + randomUUID().toString().replaceAll('-', '_');
      await sql`
        CREATE DATABASE ${sql.raw(name)} TEMPLATE acacia
      `.execute(connection);

      return {
        connection,
        getInstance: () => {
          return new Kysely<DB>({
            dialect: new PostgresDialect({
              pool: new Pool({
                connectionString: container
                  .getConnectionUri()
                  .replace('acacia', name),
              }),
            }),
          });
        },
        destroy: () => {
          return sql`
            DROP DATABASE ${sql.raw(name)}
          `.execute(connection);
        },
      };
    },
  };
}

export const testDatabaseContainer = await setupDatabaseContainer();
