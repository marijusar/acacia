'use client';

import { DialogTitle } from '../ui/dialog';
import { Input } from '../ui/input';
import { useTaskForm } from './task-form-context';

export const NameInput = () => {
  const {
    active: { name: isNameFieldActive },
    setActive,
    state: { name },
    setState,
  } = useTaskForm();

  if (!isNameFieldActive) {
    return (
      <DialogTitle
        onClick={(e) => {
          e.stopPropagation();
          setActive('name');
        }}
        className="text-xl font-semibold cursor-pointer"
      >
        <span className="text-muted-foreground">
          {name ? name : 'Edit title...'}
        </span>
      </DialogTitle>
    );
  }

  return (
    <DialogTitle>
      <div className="flex">
        <Input
          autoFocus={true}
          onClick={(e) => e.stopPropagation()}
          className="max-w-md mr-8"
          value={name}
          onChange={(e) => setState((s) => ({ ...s, name: e.target.value }))}
          placeholder="Enter title.."
        />
      </div>
    </DialogTitle>
  );
};
