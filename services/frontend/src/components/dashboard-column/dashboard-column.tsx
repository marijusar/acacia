'use client';

import { useState } from 'react';
import { TicketCard } from '../ticket-card/ticket-card';
import { Heading4 } from '../ui/headings';
import type { Issue, ProjectStatusColumn } from '@/lib/schemas/projects';

type DashboardColumnProps = ProjectStatusColumn & { issues: Issue[] };

export const DashboardColumn = ({ name, issues }: DashboardColumnProps) => {
  const [isDragTargetOver, setIsDragTargetOver] = useState(false);
  return (
    <div
      onDrop={() => setIsDragTargetOver(false)}
      onDragEnter={() => setIsDragTargetOver(true)}
      onDragLeave={() => setIsDragTargetOver(false)}
      onDragOver={(e) => e.preventDefault()}
      className="h-full min-w-64 w-full bg-card mr-8 last:mr-0 rounded-md border p-3 border-accent"
    >
      <Heading4 className="mb-6">{name}</Heading4>
      {isDragTargetOver ? <p>Target over here</p> : null}
      {issues.map((issue) => (
        <TicketCard {...issue} key={issue.id} />
      ))}
    </div>
  );
};
