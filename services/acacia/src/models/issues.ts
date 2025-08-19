import type { Kysely, Selectable } from 'kysely';
import type { DB, Issues } from '#dashboard-api/types/database.js';
import type {
  CreateIssueInput,
  UpdateIssueInput,
} from '#dashboard-api/schemas/issues.ts';

type SelectableIssue = Selectable<Issues>;

export class IssuesModel {
  database;
  constructor(database: Kysely<DB>) {
    this.database = database;
  }
  async findAll(): Promise<SelectableIssue[]> {
    return this.database
      .selectFrom('issues')
      .selectAll()
      .orderBy('created_at', 'desc')
      .execute();
  }

  async findById(id: string): Promise<SelectableIssue | undefined> {
    return this.database
      .selectFrom('issues')
      .selectAll()
      .where('id', '=', id)
      .executeTakeFirst();
  }

  async create(data: CreateIssueInput): Promise<SelectableIssue> {
    const now = new Date();

    return this.database
      .insertInto('issues')
      .values({
        ...data,
        created_at: now,
        updated_at: now,
      })
      .returningAll()
      .executeTakeFirstOrThrow();
  }

  async update(
    id: string,
    data: UpdateIssueInput
  ): Promise<SelectableIssue | undefined> {
    return this.database
      .updateTable('issues')
      .set({
        ...data,
        updated_at: new Date(),
      })
      .where('id', '=', id)
      .returningAll()
      .executeTakeFirst();
  }

  async delete(id: string): Promise<SelectableIssue | undefined> {
    return this.database
      .deleteFrom('issues')
      .where('id', '=', id)
      .returningAll()
      .executeTakeFirst();
  }
}
