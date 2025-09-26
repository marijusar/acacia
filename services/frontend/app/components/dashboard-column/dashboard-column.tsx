import { useState } from 'react';
import { TicketCard } from '../ticket-card/ticket-card';
import { Heading4 } from '../ui/headings';
import type { ProjectStatusColumn } from '~/schemas/projects';

type DashboardColumnProps = ProjectStatusColumn;

export const DashboardColumn = ({ name }: DashboardColumnProps) => {
  const [isDragTargetOver, setIsDragTargetOver] = useState(false);
  return (
    <div
      onDrop={() => setIsDragTargetOver(false)}
      onDragEnter={() => setIsDragTargetOver(true)}
      onDragLeave={() => setIsDragTargetOver(false)}
      onDragOver={(e) => e.preventDefault()}
      className="h-full min-w-64 w-full bg-card mr-8 rounded-md border p-3 border-accent"
    >
      <Heading4 className="mb-6">{name}</Heading4>
      {isDragTargetOver ? <p>Target over here</p> : null}
      <TicketCard />
    </div>
  );
};
