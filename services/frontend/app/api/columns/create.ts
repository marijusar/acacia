import { projectService } from '~/services/project-service.server';
import { createProjectColumnParams } from '~/schemas/projects';
import { redirect } from 'react-router';
import type { Route } from './+types/create';

export async function action({ request }: Route.ActionArgs) {
  const formData = await request.formData();
  const columnName = formData.get('column-name');
  const projectId = formData.get('project-id');
  const body = {
    name: columnName,
    project_id: projectId,
  };

  await projectService.createProjectColumn(
    createProjectColumnParams.parse(body)
  );

  return redirect(`/projects/board/${projectId}/settings/columns`);
}
