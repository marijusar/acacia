import { dialect } from '#dashboard-api/config/database.ts';
import { defineConfig } from 'kysely-ctl';

export default defineConfig({
  dialect,
});
