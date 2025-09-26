import { Form } from 'react-router';
import { Label } from '../ui/label';
import { Input } from '../ui/input';
import { Button } from '../ui/button';

export const ProjectForm = () => {
  return (
    <Form className="flex flex-col" method="post" action="/api/projects/create">
      <Label className="mb-4" htmlFor="project-name">
        Name
      </Label>
      <Input
        name="project-name"
        id="project-name"
        type="text"
        min="2"
        max="255"
      />

      <Button className="mt-8 ml-auto" type="submit">
        Submit
      </Button>
    </Form>
  );
};
