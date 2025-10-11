import { DashboardColumn } from '../dashboard-column/dashboard-column';
import { type Issue, type ProjectStatusColumn } from '@/lib/schemas/projects';

type DashboardColumnsProps = {
  columns: ProjectStatusColumn[];
  columnIssueMap: Record<number, Issue[]>;
};

export const DashboardColumns = ({
  columns,
  columnIssueMap,
}: DashboardColumnsProps) => {
  return (
    <div className="flex min-w-full flex-1 ml-auto mr-auto">
      {columns.map((column) => (
        <DashboardColumn
          {...column}
          key={column.id}
          issues={columnIssueMap[column.id]}
          columns={columns}
        />
      ))}
    </div>
  );
};
