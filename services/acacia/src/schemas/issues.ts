import { z } from 'zod';

export const createIssueSchema = z.object({
  name: z.string().min(1, 'Name is required').max(255, 'Name must be less than 255 characters'),
  description: z.string().optional(),
});

export const updateIssueSchema = z.object({
  name: z.string().min(1, 'Name is required').max(255, 'Name must be less than 255 characters').optional(),
  description: z.string().optional(),
});

export const issueParametersSchema = z.object({
  id: z.string().regex(/^\d+$/, 'ID must be a valid number'),
});

export type CreateIssueInput = z.infer<typeof createIssueSchema>;
export type UpdateIssueInput = z.infer<typeof updateIssueSchema>;
export type IssueParameters = z.infer<typeof issueParametersSchema>;