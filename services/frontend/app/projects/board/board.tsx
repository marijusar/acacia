import { CreateTaskDialog } from '~/components/create-task-dialog/create-task-dialog';
import { DashboardColumns } from '~/components/dashboard-columns/dashboard-columns';
import { Heading1 } from '~/components/ui/headings';
import { useDashboardContext } from '~/layouts/dashboard';

const Board = () => {
  const { projectDetails } = useDashboardContext();
  return (
    <>
      <div className="flex items-center justify-between mb-8 sticky">
        <Heading1>Sprint 1</Heading1>
        <CreateTaskDialog />
      </div>
      <DashboardColumns columns={projectDetails.columns} />
    </>
  );
};

export default Board;
