import { Label } from '../ui/label';
import { Input } from '../ui/input';
import { Button } from '../ui/button';
import { createProjectColumnAction } from '@/app/actions/projects';

type CreateColumnFormProps = {
  projectId: number;
  onSubmitComplete: () => void;
};

export const CreateColumnForm = ({
  projectId,
  onSubmitComplete,
}: CreateColumnFormProps) => {
  const createProjectColumnActionHandler = async (formData: FormData) => {
    await createProjectColumnAction(formData);
    onSubmitComplete();
  };
  return (
    <form className="flex flex-col" action={createProjectColumnActionHandler}>
      <Label className="mb-4" htmlFor="column-name">
        Name
      </Label>
      <Input
        name="column-name"
        id="column-name"
        type="text"
        min="2"
        max="255"
        required
      />

      <Input
        type="hidden"
        name="project-id"
        id="project-id"
        value={projectId}
      />

      <Button className="mt-8 ml-auto cursor-pointer" type="submit">
        Submit
      </Button>
    </form>
  );
};
