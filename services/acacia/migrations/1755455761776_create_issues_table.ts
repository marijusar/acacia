import type { Kysely } from 'kysely';

// `any` is required here since migrations should be frozen in time. alternatively, keep a "snapshot" db interface.
export async function up(database: Kysely<any>): Promise<void> {
  await database.schema
    .createTable('issues')
    .addColumn('id', 'bigserial', (c) => c.primaryKey())
    .addColumn('name', 'text')
    .addColumn('description', 'text')
    .addColumn('created_at', 'timestamptz')
    .addColumn('updated_at', 'timestamptz')
    .execute();
}

// `any` is required here since migrations should be frozen in time. alternatively, keep a "snapshot" db interface.
export async function down(database: Kysely<any>): Promise<void> {
  await database.schema.dropTable('issues').execute();
}
