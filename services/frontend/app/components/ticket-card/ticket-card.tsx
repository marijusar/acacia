import { useState } from 'react';
import { TicketCardPriority } from '../ticket-card-priority/ticket-card-priority';
import { Card, CardFooter, CardHeader, CardTitle } from '../ui/card';

export const TicketCard = () => {
  const [dragging, setDragging] = useState(false);
  const bgColor = dragging ? 'bg-primary' : 'bg-secondary';

  return (
    <Card
      draggable={true}
      className={`cursor-move p-4 rounded-md ${bgColor}`}
      onDragStart={(e) => e.dataTransfer.setData('copy', 'hello-world')}
      onDragEnd={() => setDragging(false)}
      onDrop={(e) => console.log(e)}
    >
      <CardHeader className="p-0 select-none">
        <CardTitle className="font-normal leading-6">
          Add proxy to the Zoho CRM
        </CardTitle>
      </CardHeader>
      <CardFooter className="p-0 select-none">
        <TicketCardPriority />
      </CardFooter>
    </Card>
  );
};
