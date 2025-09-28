'use client';

import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogTrigger } from '@/components/ui/dialog';
import { DialogTitle } from '@radix-ui/react-dialog';
import { MouseEvent, useState } from 'react';
import { Input } from '../ui/input';
import { Textarea } from '../ui/textarea';
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '../ui/select';
import { ProjectStatusColumn } from '@/lib/schemas/projects';
import { Label } from '../ui/label';
import { createIssueAction } from '@/app/actions/issues';

type ColumnSelectProps = {
  columns: ProjectStatusColumn[];
};

const ColumnSelect = ({ columns }: ColumnSelectProps) => {
  return (
    <div>
      <Label className="mb-2">Status</Label>
      <Select name="column_id">
        <SelectTrigger className="w-[180px] cursor-pointer">
          <SelectValue placeholder="Select status" />
        </SelectTrigger>
        <SelectContent>
          <SelectGroup>
            {columns.map(({ id, name }) => (
              <SelectItem
                key={id}
                value={id.toString()}
                className="cursor-pointer"
              >
                {name}
              </SelectItem>
            ))}
          </SelectGroup>
        </SelectContent>
      </Select>
    </div>
  );
};

type CreateTaskDialogProps = {
  columns: ProjectStatusColumn[];
};

export const CreateTaskDialog = ({ columns }: CreateTaskDialogProps) => {
  const [open, setOpen] = useState(false);
  const [formState, setFormState] = useState({ title: '', description: '' });
  const [activeFormFields, setActiveFormFields] = useState({
    title: false,
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
    <Dialog onOpenChange={() => setOpen((o) => !o)} open={open}>
      <DialogTrigger asChild>
        <Button className="ml-auto">Create task</Button>
      </DialogTrigger>
      <DialogContent
        className="flex min-w-4xl w-full min-h-96"
        onClick={onDeactivateInputs}
        aria-describedby={undefined}
      >
        <form
          action={async (formData) => {
            await createIssueAction(formData);
            setOpen(false);
          }}
          className="flex w-full"
        >
          <Input type="hidden" name="name" value={formState.title} />
          <Input
            type="hidden"
            name="description"
            value={formState.description}
          />
          <div className="flex flex-col flex-1">
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
                  className="max-w-lg w-full p-4 mr-8 min-h-24 mt-8"
                  value={formState.description}
                  placeholder="Edit description..."
                  onChange={(e) =>
                    setFormState((s) => ({ ...s, description: e.target.value }))
                  }
                />
              </div>
            ) : (
              <p
                onClick={(e) => onActivateInput(e, 'description')}
                className="max-w-md text-muted-foreground mt-8"
              >
                {formState.description
                  ? formState.description
                  : 'Edit description...'}
              </p>
            )}
          </div>
          <div className="flex flex-col ml-auto pt-12">
            <ColumnSelect columns={columns} />
            <Button type="submit" className="max-w-24 ml-auto mt-auto">
              Save
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
};
