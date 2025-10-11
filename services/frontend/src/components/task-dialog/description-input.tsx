'use client';

import { Textarea } from '../ui/textarea';
import { useTaskForm } from './task-form-context';

export const DescriptionInput = () => {
  const { active, state, setActive, setState } = useTaskForm();
  if (active.description) {
    return (
      <div className="flex">
        <Textarea
          autoFocus={true}
          name="description"
          onClick={(e) => e.stopPropagation()}
          className="max-w-lg w-full p-4 mr-8 min-h-24 mt-8"
          value={state.description}
          placeholder="Edit description..."
          onChange={(e) => {
            setState((s) => ({ ...s, description: e.target.value }));
          }}
        />
      </div>
    );
  }
  return (
    <p
      onClick={(e) => {
        e.stopPropagation();
        setActive('description');
      }}
      className="max-w-md text-muted-foreground mt-8 block cursor-pointer"
    >
      {state.description ? state.description : 'Edit description...'}
    </p>
  );
};
