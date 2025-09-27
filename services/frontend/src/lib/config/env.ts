import { envSchema } from '@/lib/schemas/env';

const unsafeEnv = {
  ACACIA_API_URL: process.env.ACACIA_API_URL,
};

export const env = envSchema.parse(unsafeEnv);