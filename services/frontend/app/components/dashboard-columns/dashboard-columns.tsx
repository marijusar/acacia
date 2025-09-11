import { DashboardColumn } from '../dashboard-column/dashboard-column';

export const DashboardColumns = () => {
  return (
    <div className="flex w-full flex-1">
      <DashboardColumn />
      <DashboardColumn />
    </div>
  );
};
