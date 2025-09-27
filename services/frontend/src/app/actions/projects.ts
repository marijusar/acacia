'use server';

import { projectService } from '@/lib/services/project-service';
import {
  createProjectParams,
  createProjectColumnParams,
  deleteProjectColumnParams,
} from '@/lib/schemas/projects';
import { redirect } from 'next/navigation';
import { revalidatePath } from 'next/cache';

export async function createProjectAction(formData: FormData) {
  const projectName = formData.get('project-name');
  const body = {
    name: projectName,
  };

  const project = await projectService.createProject(
    createProjectParams.parse(body)
  );

  // Revalidate the projects list cache
  revalidatePath('/projects');

  // Redirect to the new project board
  redirect(`/projects/${project.id}/board`);
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

