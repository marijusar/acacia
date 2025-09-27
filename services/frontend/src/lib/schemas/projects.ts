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

export type Project = z.infer<typeof project>;

export const projectDetailsResponse = z.object({
  ...project.shape,
  columns: z.array(projectStatusColumn),
});

export const projectsResponse = z.array(project);

export const createProjectParams = z.object({ name: z.string() });

export const createProjectResponse = z.object({
  id: z.number(),
  name: z.string(),
  created_at: z.string(),
  updated_at: z.string(),
});

export const projectDashboardRouteArguments = z.object({ id: z.string() });

export const projectColumnsResponse = z.array(
  z.object({
    id: z.number(),
    name: z.string(),
    position_index: z.number(),
    created_at: z.string(),
    updated_at: z.string(),
  })
);

export const createProjectColumnParams = z.object({
  project_id: z.string().transform((v) => Number(v)),
  name: z.string(),
});

export const deleteProjectColumnParams = z.object({
  'project-column-id': z.string().transform((v) => Number(v)),
});