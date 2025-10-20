'use client';

import { DragEvent, useState } from 'react';
import { TicketCardPriority } from '../ticket-card-priority/ticket-card-priority';
import { Card, CardFooter, CardHeader, CardTitle } from '../ui/card';
import { issue, Issue } from '@/lib/schemas/projects';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';

type TicketCardProps = Issue;

export const TicketCard = ({
  id,
  name,
  description,
  column_id,
  description_serialized,
  ...rest
}: TicketCardProps) => {
  const params = useSearchParams();
  const pathname = usePathname();
  const router = useRouter();
  const [dragging, setDragging] = useState(false);
  const bgColor = dragging ? 'bg-primary' : 'bg-secondary';
  const onDragStart = (e: DragEvent<HTMLDivElement>) => {
    e.dataTransfer.setData(
      'text/plain',
      JSON.stringify(
        issue.parse({
          id,
          name,
          description,
          description_serialized,
          column_id,
        })
      )
    );
  };

  const handleOpenCard = () => {
    const searchParams = new URLSearchParams(params);
    searchParams.append('open_issue_id', id.toString());
    router.push(`${pathname}?${searchParams.toString()}`);
  };

  return (
    <Card
      onClick={handleOpenCard}
      draggable={true}
      className={`cursor-move p-4 rounded-md ${bgColor} mb-4`}
      onDragStart={onDragStart}
      onDragEnd={() => setDragging(false)}
      onDrop={(e) => console.log(e)}
      {...rest}
    >
      <CardHeader className="p-0 select-none">
        <CardTitle className="font-normal leading-6">{name}</CardTitle>
      </CardHeader>
      <CardFooter className="p-0 select-none">
        <TicketCardPriority />
      </CardFooter>
    </Card>
  );
};
