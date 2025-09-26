import { DashboardColumns } from '~/components/dashboard-columns/dashboard-columns';
import { Heading1 } from '~/components/ui/headings';
import { useDashboardContext } from '~/layouts/dashboard';

const Board = () => {
  const { projectDetails } = useDashboardContext();
  return (
    <>
      <Heading1 className="mb-8">Sprint 1</Heading1>
      <DashboardColumns columns={projectDetails.columns} />
    </>
  );
};

export default Board;
