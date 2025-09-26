import { projectService } from '~/services/project-service.server';

export async function action({ request }: { request: Request }) {
  const formData = await request.formData();
  const columnId = formData.get('project-column-id');

  await projectService.deleteProjectColumn(Number(columnId));

  return null;
}
