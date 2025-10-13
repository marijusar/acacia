import z from 'zod';

export const teamResponse = z.object({
  id: z.number(),
  name: z.string(),
  description: z.string().optional().nullable(),
  created_at: z.string(),
  updated_at: z.string(),
});

export type TeamResponse = z.infer<typeof teamResponse>;

export const userTeamsResponse = z.array(teamResponse);

export type UserTeamsResponse = z.infer<typeof userTeamsResponse>;

export const createTeamInput = z.object({
  name: z
    .string()
    .min(1, 'Team name is required')
    .max(255, 'Team name must be less than 255 characters'),
});

export type CreateTeamInput = z.infer<typeof createTeamInput>;

export const createTeamResponse = teamResponse;

export type CreateTeamResponse = z.infer<typeof createTeamResponse>;
