'use client';

import {
  createContext,
  Dispatch,
  ReactNode,
  SetStateAction,
  useContext,
  useState,
} from 'react';
import { Dialog, DialogContent, DialogTrigger } from '../ui/dialog';
import { Issue } from '@/lib/schemas/projects';
import { EditorState, SerializedEditorState, $getRoot } from 'lexical';

const activeFormFieldsDefaultState = { name: false, description: false };

type CreateDialogFormState = {
  name: string;
  description: EditorState | null;
  column_id: string | undefined;
};

const initialState = {
  name: '',
  description: null,
  column_id: undefined,
} satisfies CreateDialogFormState;

const TaskFormContext = createContext<{
  active: typeof activeFormFieldsDefaultState;
  state: CreateDialogFormState;
  setActive: (field: keyof CreateDialogFormState) => void;
  setState: Dispatch<SetStateAction<CreateDialogFormState>>;
  initialSerializedState?: SerializedEditorState | null;
} | null>(null);

type TaskFormProps = {
  action: (formData: FormData) => Promise<void>;
  trigger: ReactNode;
  children: ReactNode;
  issue?: Issue;
};

export const TaskFormProvider = ({
  action,
  trigger,
  children,
  issue,
}: TaskFormProps) => {
  const [open, setOpen] = useState(false);
  console.log(issue);

  const [formState, setFormState] = useState<CreateDialogFormState>(
    issue
      ? {
          ...issue,
          column_id: issue?.column_id.toString(),
          description: initialState.description,
        }
      : initialState
  );

  // Parse serialized state if editing an existing issue
  const initialSerializedState = issue?.description_serialized
    ? (JSON.parse(issue.description_serialized) as SerializedEditorState)
    : null;

  const setFieldActive = (field: keyof CreateDialogFormState) => {
    setActiveFormFields((fields) => ({
      ...fields,
      [field]: true,
    }));
  };

  const [activeFormFields, setActiveFormFields] = useState(
    activeFormFieldsDefaultState
  );

  const onDeactivateInputs = () => {
    setActiveFormFields(activeFormFieldsDefaultState);
  };

  return (
    <Dialog onOpenChange={() => setOpen((o) => !o)} open={open}>
      <DialogTrigger asChild>{trigger}</DialogTrigger>
      <DialogContent
        className="flex min-w-4xl w-full min-h-96"
        onClick={onDeactivateInputs}
        aria-describedby={undefined}
      >
        <form
          action={async (formData) => {
            // Extract plain text and serialized state from EditorState before submitting
            if (formState.description) {
              const plainText = formState.description.read(() =>
                $getRoot().getTextContent()
              );
              const serializedState = JSON.stringify(
                formState.description.toJSON()
              );

              formData.set('description', plainText);
              formData.set('description_serialized', serializedState);
            }

            await action(formData);
            setOpen(false);
            setFormState(initialState);
          }}
          className="flex w-full"
        >
          <input type="hidden" name="id" value={issue?.id} />
          <input type="hidden" name="name" value={formState.name} />
          <input type="hidden" name="description" value="" />
          <input type="hidden" name="description_serialized" value="" />
          <TaskFormContext.Provider
            value={{
              active: activeFormFields,
              setActive: setFieldActive,
              state: formState,
              setState: setFormState,
              initialSerializedState,
            }}
          >
            {children}
          </TaskFormContext.Provider>
        </form>
      </DialogContent>
    </Dialog>
  );
};

export const useTaskForm = () => {
  const ctx = useContext(TaskFormContext);

  if (!ctx) {
    throw new Error('useTaskForm must be used within TaskFormProvider');
  }

  return ctx;
};
