import { projectService } from '@/lib/services/project-service';
import { DashboardColumns } from '@/components/dashboard-columns/dashboard-columns';
import { Heading1 } from '@/components/ui/headings';
import { Button } from '@/components/ui/button';
import { CreateTaskDialog } from '@/components/task-dialog';

interface BoardPageProps {
  params: Promise<{ id: string }>;
}

export default async function BoardPage({ params }: BoardPageProps) {
  const { id } = await params;

  const projectDetails = await projectService.getProjectDetails(id);

  if (!projectDetails) {
    return <Heading1> Not implemented </Heading1>;
  }

  return (
    <>
      <div className="flex items-center justify-between mb-8 sticky left-0">
        <Heading1>Sprint 1</Heading1>
        <CreateTaskDialog columns={projectDetails.columns}>
          <Button className="ml-auto">Create task</Button>
        </CreateTaskDialog>
      </div>
      <DashboardColumns
        columns={projectDetails.columns}
        columnIssueMap={projectDetails.issues}
      />
    </>
  );
}
