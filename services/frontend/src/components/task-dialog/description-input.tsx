'use client';

import { twMerge } from 'tailwind-merge';
import { Editor } from '../blocks/editor-00/editor';
import { useTaskForm } from './task-form-context';
import { $getRoot } from 'lexical';

export const DescriptionInput = () => {
  const { setActive, setState } = useTaskForm();
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
        onSerializedChange={(v) =>
          setState((state) => ({ ...state, description: v }))
        }
        onChange={(s) => s.read(() => $getRoot().getTextContent())}
      />
    </div>
  );
};
