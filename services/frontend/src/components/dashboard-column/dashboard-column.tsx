'use client';

import { DragEvent } from 'react';
import { TicketCard } from '../ticket-card/ticket-card';
import { Heading4 } from '../ui/headings';
import {
  issue,
  type Issue,
  type ProjectStatusColumn,
} from '@/lib/schemas/projects';
import { updateIssue } from '@/app/actions/issues';
import { UpdateTaskDialog } from '../task-dialog';

type DashboardColumnProps = ProjectStatusColumn & {
  issues: Issue[];
};

export const DashboardColumn = ({
  name,
  issues = [],
  id,
}: DashboardColumnProps) => {
  const onDrop = async (e: DragEvent<HTMLDivElement>) => {
    const droppedIssue = issue.parse(
      JSON.parse(e.dataTransfer.getData('text/plain'))
    );

    await updateIssue({ ...droppedIssue, column_id: id });
  };
  return (
    <div
      onDrop={onDrop}
      onDragOver={(e) => e.preventDefault()}
      className="h-full min-w-64 max-w-64 w-full bg-card mr-8 last:mr-0 rounded-md border p-3 border-accent"
    >
      <Heading4 className="mb-6">{name}</Heading4>
      {issues.map((issue) => (
        <TicketCard key={issue.id} {...issue} />
      ))}
    </div>
  );
};
