'use client';

import { useState } from 'react';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { PlusIcon } from 'lucide-react';
import { CreateColumnForm } from '../create-column-form/create-column-form';

type CreateColumnDialogProps = {
  projectId: number;
};

export const CreateColumnDialog = ({ projectId }: CreateColumnDialogProps) => {
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const onToggleCreateColumnModal = () => {
    setIsDialogOpen((o) => !o);
  };
  return (
    <>
      <PlusIcon
        onClick={onToggleCreateColumnModal}
        className="bg-sidebar-primary rounded-sm ml-auto cursor-pointer"
      />
      <Dialog onOpenChange={onToggleCreateColumnModal} open={isDialogOpen}>
        <DialogContent aria-describedby={undefined}>
          <DialogHeader>
            <DialogTitle>Create column</DialogTitle>
          </DialogHeader>

          <CreateColumnForm
            onSubmitComplete={() => setIsDialogOpen(false)}
            projectId={projectId}
          />
        </DialogContent>
      </Dialog>
    </>
  );
};
