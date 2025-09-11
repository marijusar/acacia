import { DashboardColumn } from '../dashboard-column/dashboard-column';
import type { ProjectStatusColumn } from '../../schemas/projects';

type DashboardColumnsProps = {
  columns: ProjectStatusColumn[];
};

export const DashboardColumns = ({ columns }: DashboardColumnsProps) => {
  return (
    <div className="flex w-full flex-1 ">
      {columns.map((column) => (
        <DashboardColumn {...column} key={column.id} />
      ))}
    </div>
  );
};
