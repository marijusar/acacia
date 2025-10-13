'use server';

import { teamService } from '@/lib/services/team-service';
import { createTeamInput } from '@/lib/schemas/teams';
import { redirect } from 'next/navigation';

export type CreateTeamFormState = {
  error: string | null;
};

export async function createTeamAction(
  _: CreateTeamFormState,
  formData: FormData
): Promise<CreateTeamFormState> {
  try {
    const data = {
      name: formData.get('name'),
    };

    const validatedData = createTeamInput.parse(data);

    await teamService.createTeam(validatedData);
  } catch (error) {
    if (error instanceof Error) {
      return {
        error: error.message,
      };
    }
    return {
      error: 'An unexpected error has occurred',
    };
  }

  redirect('/projects');
}
