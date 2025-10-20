'use server';

import { projectService } from '@/lib/services/project-service';
import { logger } from '@/lib/config/logger';
import {
  createProjectInput,
  createProjectColumnParams,
  deleteProjectColumnParams,
  type CreateProjectResponse,
} from '@/lib/schemas/projects';
import { redirect } from 'next/navigation';
import { revalidatePath } from 'next/cache';

export type CreateProjectFormState = {
  error: string | null;
};

const initialState: CreateProjectFormState = {
  error: null,
};

export async function createProjectAction(
  _: CreateProjectFormState,
  formData: FormData
): Promise<CreateProjectFormState> {
  let project: CreateProjectResponse | null = null;

  try {
    const data = {
      name: formData.get('name'),
    };

    const teamIdStr = formData.get('team_id');

    if (!teamIdStr) {
      return {
        error: 'Please select a team',
      };
    }

    const validatedData = createProjectInput.parse(data);

    project = await projectService.createProject({
      name: validatedData.name,
      team_id: parseInt(teamIdStr.toString()),
    });

    // Revalidate the projects list cache
    revalidatePath('/projects');
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

  // Redirect outside try/catch so Next.js can throw properly
  if (project?.id) {
    redirect(`/projects/${project.id}/board`);
  }

  logger.error(
    '[CREATE_PROJECT_ACTION] Project was created but id is missing',
    { project }
  );
  return {
    error: 'An unexpected error has occurred',
  };
}

export async function createProjectColumnAction(formData: FormData) {
  const body = {
    name: formData.get('column-name'),
    project_id: formData.get('project-id'),
  };

  const validatedData = createProjectColumnParams.parse(body);
  await projectService.createProjectColumn(validatedData);

  // Revalidate the project details cache
  revalidatePath(`/projects/${validatedData.project_id}`);
}

export async function deleteProjectColumnAction(formData: FormData) {
  const body = Object.fromEntries(formData);
  const validatedData = deleteProjectColumnParams.parse(body);

  await projectService.deleteProjectColumn(validatedData['project-column-id']);

  revalidatePath(`/projects/`, 'layout');
}
