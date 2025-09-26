import { projectService } from '~/services/project-service.server';
import type { Route } from './+types/create';
import { createProjectParams } from '~/schemas/projects';
import { redirect } from 'react-router';

export async function action({ request }: Route.ActionArgs) {
  const formData = await request.formData();
  const projectName = formData.get('project-name');
  const body = {
    name: projectName,
  };

  const project = await projectService.createProject(
    createProjectParams.parse(body)
  );

  return redirect(`/projects/board/${project.id}`);
}
