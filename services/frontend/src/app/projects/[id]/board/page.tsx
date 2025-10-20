import { projectService } from '@/lib/services/project-service';
import { DashboardColumns } from '@/components/dashboard-columns/dashboard-columns';
import { Heading1 } from '@/components/ui/headings';
import { CreateTaskDialog, UpdateTaskDialog } from '@/components/task-dialog';
import { CreateTaskButton } from '@/components/create-task-button/create-task-button';

interface BoardPageProps {
  params: Promise<{ id: string }>;
  searchParams: Promise<{ open_issue_id: string }>;
}

export default async function BoardPage({
  params,
  searchParams,
}: BoardPageProps) {
  const { open_issue_id } = await searchParams;
  const { id } = await params;

  const projectDetails = await projectService.getProjectDetails(id);

  if (!projectDetails) {
    return <Heading1> Not implemented </Heading1>;
  }

  return (
    <>
      <div className="flex items-center justify-between mb-8 sticky left-0">
        <Heading1>Sprint 1</Heading1>
        <CreateTaskButton />
      </div>
      <DashboardColumns
        columns={projectDetails.columns}
        columnIssueMap={projectDetails.issues}
      />

      {open_issue_id && open_issue_id !== 'new' && (
        <UpdateTaskDialog
          columns={projectDetails.columns}
          issueId={parseInt(open_issue_id)}
        />
      )}

      {open_issue_id && open_issue_id === 'new' && (
        <CreateTaskDialog columns={projectDetails.columns} />
      )}
    </>
  );
}
