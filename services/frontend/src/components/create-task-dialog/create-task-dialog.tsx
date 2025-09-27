'use client';

import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogTrigger } from '@/components/ui/dialog';
import { DialogTitle } from '@radix-ui/react-dialog';
import { MouseEvent, useState } from 'react';
import { Input } from '../ui/input';
import { Textarea } from '../ui/textarea';

export const CreateTaskDialog = () => {
  const [formState, setFormState] = useState({ title: '', description: '' });
  const [activeFormFields, setActiveFormFields] = useState({
    title: true,
    description: false,
  });

  const onDeactivateInputs = () => {
    setActiveFormFields({ title: false, description: false });
  };

  const onActivateInput = (
    e: MouseEvent,
    type: keyof typeof activeFormFields
  ) => {
    e.stopPropagation();
    setActiveFormFields((f) => ({ ...f, [type]: true }));
  };

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button>Create task</Button>
      </DialogTrigger>
      <DialogContent
        className="min-w-4xl w-full"
        onClick={onDeactivateInputs}
        aria-describedby={undefined}
      >
        {activeFormFields.title ? (
          <DialogTitle>
            <div className="flex">
              <Input
                onClick={(e) => e.stopPropagation()}
                className="max-w-md mr-8"
                value={formState.title}
                onChange={(e) =>
                  setFormState((s) => ({ ...s, title: e.target.value }))
                }
                placeholder="Enter title.."
              />
            </div>
          </DialogTitle>
        ) : (
          <DialogTitle
            onClick={(e) => onActivateInput(e, 'title')}
            className="text-xl font-semibold"
          >
            {formState.title ? (
              formState.title
            ) : (
              <span className="text-muted-foreground">Edit title...</span>
            )}
          </DialogTitle>
        )}
        {activeFormFields.description ? (
          <div className="flex">
            <Textarea
              onClick={(e) => e.stopPropagation()}
              className="max-w-lg w-full p-4 mr-8 min-h-24 "
              value={formState.description}
              onChange={(e) =>
                setFormState((s) => ({ ...s, description: e.target.value }))
              }
            />
          </div>
        ) : (
          <p
            onClick={(e) => onActivateInput(e, 'description')}
            className="max-w-md text-muted-foreground"
          >
            {formState.description
              ? formState.description
              : 'Edit description...'}
          </p>
        )}
        <Button className="max-w-24 ml-auto">Save</Button>
      </DialogContent>
    </Dialog>
  );
};
