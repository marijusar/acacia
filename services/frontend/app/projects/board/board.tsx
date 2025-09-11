import { DashboardColumns } from '~/components/dashboard-columns/dashboard-columns';
import { AppSidebar } from '~/components/sidebar/sidebar';
import { Heading1 } from '~/components/ui/headings';
import { Input } from '~/components/ui/input';

const Board = () => {
  return (
    <div className="flex">
      <AppSidebar />
      <div className="w-full flex flex-col  ">
        <div className="h-16 w-full bg-secondary flex items-center justify-center">
          <Input placeholder="Search.." className="max-w-100" />
        </div>
        <div className="w-24"> </div>
        <div className="pt-4 pb-4 pr-8 pl-8 h-full flex flex-col">
          <Heading1 className="mb-8">Sprint 1</Heading1>
          <DashboardColumns />
        </div>
      </div>
    </div>
  );
};

export default Board;
