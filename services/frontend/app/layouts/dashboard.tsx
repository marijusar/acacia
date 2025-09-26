import { projectService } from '~/services/project-service.server';
import type { Route } from '../+types/root';
import { Outlet, useLoaderData, useOutletContext } from 'react-router';
import { AppSidebar } from '~/components/sidebar/sidebar';
import { Input } from '~/components/ui/input';
import { projectDashboardRouteArguments } from '~/schemas/projects';

export async function loader({ params }: Route.LoaderArgs) {
  const safeParams = projectDashboardRouteArguments.parse(params);
  const [projectDetails, projects] = await Promise.all([
    projectService.getProjectDetails(safeParams.id),
    projectService.getProjects(),
  ]);
  return { projectDetails, projects };
}

type AwaitedLoader = Awaited<ReturnType<typeof loader>>;

const DashboardLayout = () => {
  const { projectDetails, projects } = useLoaderData<typeof loader>();
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
          <Outlet context={{ projectDetails, projects }} />
        </div>
      </div>
    </div>
  );
};

export const useDashboardContext = () => {
  const ctx = useOutletContext<AwaitedLoader>();
  return ctx;
};

export default DashboardLayout;
