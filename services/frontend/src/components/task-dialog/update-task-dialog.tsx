import { updateIssueForm } from '@/app/actions/issues';
import { TaskFormProvider } from './task-form-context';
import { NameInput } from './name-input';
import { DescriptionInput } from './description-input';
import { SelectColumn } from './select-column';
import { Button } from '../ui/button';
import { Issue, ProjectStatusColumn } from '@/lib/schemas/projects';
import { ReactNode } from 'react';

type CreateTaskDialogProps = {
  children: ReactNode;
  columns: ProjectStatusColumn[];
  issue: Issue;
};

export const UpdateTaskDialog = ({
  children,
  columns,
  issue,
}: CreateTaskDialogProps) => {
  return (
    <TaskFormProvider issue={issue} trigger={children} action={updateIssueForm}>
      <div className="w-full flex">
        <div className="flex flex-col flex-1">
          <NameInput />
          <DescriptionInput />
        </div>

        <div className="flex flex-col ml-auto pt-12">
          <div className="flex flex-col ml-auto flex-1">
            <SelectColumn columns={columns} />
            <Button type="submit" className="max-w-24 ml-auto mt-auto">
              Save
            </Button>
          </div>
        </div>
      </div>
    </TaskFormProvider>
  );
};
