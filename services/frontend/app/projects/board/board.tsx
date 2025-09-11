import { DashboardColumns } from '~/components/dashboard-columns/dashboard-columns';
import { AppSidebar } from '~/components/sidebar/sidebar';
import { Heading1 } from '~/components/ui/headings';
import { Input } from '~/components/ui/input';
import { projectService } from '~/services/project-service';
import type { Route } from './+types/board';

export async function loader({ params }: Route.LoaderArgs) {
  const projectDetails = projectService.getProjectDetails(params.id);

  return projectDetails;
}

const Board = ({ loaderData: { columns, name } }: Route.ComponentProps) => {
  return (
    <div className="flex">
      <AppSidebar name={name} />
      <div className="w-full flex flex-col  ">
        <div className="h-16 w-full bg-secondary flex items-center justify-center">
          <Input placeholder="Search.." className="max-w-100" />
        </div>
        <div className="w-24"> </div>
        <div className="pt-4 pb-4 pr-8 pl-8 h-full flex flex-col bg-background">
          <Heading1 className="mb-8">Sprint 1</Heading1>
          <DashboardColumns columns={columns} />
        </div>
      </div>
    </div>
  );
};

export default Board;
