import { useState } from 'react';
import { TicketCard } from '../ticket-card/ticket-card';
import { Heading4 } from '../ui/headings';

export const DashboardColumn = () => {
  const [isDragTargetOver, setIsDragTargetOver] = useState(false);
  return (
    <div
      onDrop={(e) => setIsDragTargetOver(false)}
      onDragEnter={(e) => setIsDragTargetOver(true)}
      onDragLeave={(e) => setIsDragTargetOver(false)}
      onDragOver={(e) => e.preventDefault()}
      className="h-full max-w-64 w-full bg-card mr-8 rounded-md border p-3 border-accent"
    >
      <Heading4 className="mb-6">To Do</Heading4>
      {isDragTargetOver ? <p>Target over here</p> : null}
      <TicketCard />
    </div>
  );
};
