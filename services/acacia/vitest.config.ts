import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    hookTimeout: 60_000,
    testTimeout: 30_000,
    globalSetup: './tests/setup/index.ts',
  },
});
