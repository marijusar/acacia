import { DashboardColumn } from '../dashboard-column/dashboard-column';
import type { ProjectStatusColumn } from '@/lib/schemas/projects';

type DashboardColumnsProps = {
  columns: ProjectStatusColumn[];
};

export const DashboardColumns = ({ columns }: DashboardColumnsProps) => {
  return (
    <div className="flex min-w-full flex-1 ml-auto mr-auto">
      {columns.map((column) => (
        <DashboardColumn {...column} key={column.id} />
      ))}
    </div>
  );
};
