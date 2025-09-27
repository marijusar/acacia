import { Label } from '../ui/label';
import { Input } from '../ui/input';
import { Button } from '../ui/button';
import { createProjectAction } from '@/app/actions/projects';

export const ProjectForm = () => {
  return (
    <form className="flex flex-col" action={createProjectAction}>
      <Label className="mb-4" htmlFor="project-name">
        Name
      </Label>
      <Input
        name="project-name"
        id="project-name"
        type="text"
        min="2"
        max="255"
        required
      />

      <Button className="mt-8 ml-auto" type="submit">
        Submit
      </Button>
    </form>
  );
};
