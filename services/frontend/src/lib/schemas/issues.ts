import z from 'zod';

export const createIssueRequestBody = z.object({
  name: z.string(),
  description: z.string(),
  description_serialized: z.string().optional(),
  column_id: z
    .string()
    .transform((s) => parseInt(s))
    .pipe(z.number()),
});
