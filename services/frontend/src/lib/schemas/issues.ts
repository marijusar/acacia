import z from 'zod';

export const createIssueRequestBody = z.object({
  name: z.string(),
  description: z.string(),
  column_id: z
    .string()
    .transform((s) => parseInt(s))
    .pipe(z.number()),
});
