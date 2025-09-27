import { projectService } from '@/lib/services/project-service';
import { redirect } from 'next/navigation';

export default async function HomePage() {
  const projects = await projectService.getProjects();

  if (projects.length > 0) {
    redirect(`/projects/${projects[0].id}/board`);
  }

  return (
    <div className="flex min-h-screen items-center justify-center">
      <div className="text-center">
        <h1 className="text-2xl font-bold mb-4">No Projects Found</h1>
        <p className="text-muted-foreground">
          Create your first project to get started.
        </p>
      </div>
    </div>
  );
}
