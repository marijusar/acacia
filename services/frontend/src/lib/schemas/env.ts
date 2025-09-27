import { z } from 'zod';

export const envSchema = z.object({
  ACACIA_API_URL: z.string(),
});