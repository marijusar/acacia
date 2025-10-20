import { updateIssueForm } from '@/app/actions/issues';
import { TaskFormProvider } from './task-form-context';
import { NameInput } from './name-input';
import { DescriptionInput } from './description-input';
import { SelectColumn } from './select-column';
import { Button } from '../ui/button';
import { ProjectStatusColumn } from '@/lib/schemas/projects';
import { issuesService } from '@/lib/services/issues-service';

type CreateTaskDialogProps = {
  columns: ProjectStatusColumn[];
  issueId: number;
};

export const UpdateTaskDialog = async ({
  columns,
  issueId,
}: CreateTaskDialogProps) => {
  const issue = await issuesService.getIssueById(issueId);
  return (
    <TaskFormProvider issue={issue} action={updateIssueForm}>
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
