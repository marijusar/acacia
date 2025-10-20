import { createIssueAction } from '@/app/actions/issues';
import { TaskFormProvider } from './task-form-context';
import { NameInput } from './name-input';
import { DescriptionInput } from './description-input';
import { SelectColumn } from './select-column';
import { Button } from '../ui/button';
import { ProjectStatusColumn } from '@/lib/schemas/projects';

type CreateTaskDialogProps = {
  columns: ProjectStatusColumn[];
};

export const CreateTaskDialog = async ({ columns }: CreateTaskDialogProps) => {
  return (
    <TaskFormProvider action={createIssueAction}>
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
