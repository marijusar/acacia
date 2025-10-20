'use client';

import { twMerge } from 'tailwind-merge';
import { Editor } from '../blocks/editor-00/editor';
import { useTaskForm } from './task-form-context';

export const DescriptionInput = () => {
  const { setActive, setState, initialSerializedState } = useTaskForm();
  console.log(initialSerializedState);
  return (
    <div
      onClick={(e) => {
        e.stopPropagation();
        setActive('description');
      }}
      className="flex w-full"
    >
      <Editor
        className={twMerge(
          'w-full min-h-36 mt-8 mr-4 focus-within:border-primary',
          'border shadow'
        )}
        editorSerializedState={initialSerializedState || undefined}
        onChange={(editorState) =>
          setState((state) => ({ ...state, description: editorState }))
        }
      />
    </div>
  );
};
