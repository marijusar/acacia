import z from 'zod';

export const createTeamLLMAPIKeyInput = z.object({
  provider: z.string().min(1, 'Provider is required').max(50),
  api_key: z.string().min(1, 'API key is required'),
});

export type CreateTeamLLMAPIKeyInput = z.infer<
  typeof createTeamLLMAPIKeyInput
>;

export const saveAPIKeyFormSchema = createTeamLLMAPIKeyInput.extend({
  team_id: z.string().transform((val) => parseInt(val, 10)),
});

export type SaveAPIKeyFormData = z.infer<typeof saveAPIKeyFormSchema>;

export const teamLLMAPIKeyStatusResponse = z.object({
  id: z.number(),
  provider: z.string(),
  is_active: z.boolean(),
  created_at: z.string(),
  updated_at: z.string(),
  last_used_at: z.string().optional().nullable(),
});

export type TeamLLMAPIKeyStatusResponse = z.infer<
  typeof teamLLMAPIKeyStatusResponse
>;

export const teamLLMAPIKeysListResponse = z.array(teamLLMAPIKeyStatusResponse);

export type TeamLLMAPIKeysListResponse = z.infer<
  typeof teamLLMAPIKeysListResponse
>;
