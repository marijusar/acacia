import { projectService } from '@/lib/services/project-service';
import { CreateTaskDialog } from '@/components/create-task-dialog/create-task-dialog';
import { DashboardColumns } from '@/components/dashboard-columns/dashboard-columns';
import { Heading1 } from '@/components/ui/headings';

interface BoardPageProps {
  params: Promise<{ id: string }>;
}

export default async function BoardPage({ params }: BoardPageProps) {
  const { id } = await params;

  const projectDetails = await projectService.getProjectDetails(id);

  return (
    <>
      <div className="flex items-center justify-between mb-8 sticky">
        <Heading1>Sprint 1</Heading1>
        <CreateTaskDialog columns={projectDetails.columns} />
      </div>
      <DashboardColumns
        columns={projectDetails.columns}
        columnIssueMap={projectDetails.issues}
      />
    </>
  );
}
