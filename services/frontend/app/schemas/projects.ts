import z from 'zod';

export const projectStatusColumn = z.object({
  id: z.number(),
  name: z.string(),
  project_id: z.number(),
  position_index: z.number(),
  created_at: z.string(),
  updated_at: z.string(),
});

export type ProjectStatusColumn = z.infer<typeof projectStatusColumn>;

export const project = z.object({
  id: z.number(),
  name: z.string(),
  created_at: z.string(),
  updated_at: z.string(),
});

export const projectDetailsResponse = z.object({
  ...project.shape,
  columns: z.array(projectStatusColumn),
});
