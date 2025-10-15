'use server';

import { teamLLMAPIKeysService } from '@/lib/services/team-llm-api-keys-service';
import { saveAPIKeyFormSchema } from '@/lib/schemas/team-llm-api-keys';
import { revalidatePath } from 'next/cache';

export type ActionResult = {
  success: boolean;
  error?: string;
};

export async function saveAPIKeyAction(
  prevState: ActionResult,
  formData: FormData
): Promise<ActionResult> {
  const rawData = {
    team_id: formData.get('team_id'),
    provider: formData.get('provider'),
    api_key: formData.get('api_key'),
  };

  const parsed = saveAPIKeyFormSchema.safeParse(rawData);
  if (!parsed.success) {
    return {
      success: false,
      error: parsed.error.errors[0]?.message || 'Validation failed',
    };
  }

  try {
    await teamLLMAPIKeysService.createOrUpdateAPIKey(parsed.data.team_id, {
      provider: parsed.data.provider,
      api_key: parsed.data.api_key,
    });

    revalidatePath(`/teams/${parsed.data.team_id}/settings/api-keys`);

    return { success: true };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Failed to save API key',
    };
  }
}

export async function deleteAPIKeyAction(
  teamId: number,
  keyId: number
): Promise<ActionResult> {
  try {
    await teamLLMAPIKeysService.deleteAPIKey(teamId, keyId);

    revalidatePath(`/teams/${teamId}/settings/api-keys`);

    return { success: true };
  } catch (error) {
    return {
      success: false,
      error:
        error instanceof Error ? error.message : 'Failed to delete API key',
    };
  }
}
