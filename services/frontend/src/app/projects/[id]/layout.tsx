import { projectService } from '@/lib/services/project-service';
import { AppSidebar } from '@/components/sidebar/sidebar';
import { Input } from '@/components/ui/input';

interface ProjectLayoutProps {
  children: React.ReactNode;
  params: Promise<{ id: string }>;
}

export default async function ProjectLayout({
  children,
  params,
}: ProjectLayoutProps) {
  const { id } = await params;

  const [projectDetails, projects] = await Promise.all([
    projectService.getProjectDetails(id),
    projectService.getProjects(),
  ]);

  return (
    <div className="flex h-full overflow-y-hidden">
      <AppSidebar
        projects={projects}
        projectName={projectDetails.name}
        projectId={projectDetails.id}
      />
      <div className="flex flex-col flex-1 overflow-hidden">
        <div className="h-16 w-full bg-secondary flex items-center justify-center sticky">
          <Input placeholder="Search.." className="max-w-100" />
        </div>
        <div className="w-24"> </div>
        <div className="flex pt-4 pb-4 pr-8 pl-8 flex-1 flex-col bg-background overflow-auto">
          {children}
        </div>
      </div>
    </div>
  );
}
