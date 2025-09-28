'use client';

import { useState } from 'react';
import { TicketCardPriority } from '../ticket-card-priority/ticket-card-priority';
import { Card, CardFooter, CardHeader, CardTitle } from '../ui/card';
import { Issue } from '@/lib/schemas/projects';

type TicketCardProps = Issue;

export const TicketCard = ({ name }: TicketCardProps) => {
  const [dragging, setDragging] = useState(false);
  const bgColor = dragging ? 'bg-primary' : 'bg-secondary';

  return (
    <Card
      draggable={true}
      className={`cursor-move p-4 rounded-md ${bgColor} mb-4`}
      onDragStart={(e) => e.dataTransfer.setData('copy', 'hello-world')}
      onDragEnd={() => setDragging(false)}
      onDrop={(e) => console.log(e)}
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
