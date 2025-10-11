import z from 'zod';

export const dropEventData = z.object({
  columnId: z.number(),
  issueId: z.number(),
});
