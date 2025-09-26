import { Form } from 'react-router';
import { Label } from '../ui/label';
import { Input } from '../ui/input';
import { Button } from '../ui/button';

type CreateColumnFormProps = {
  projectId: number;
};

export const CreateColumnForm = ({ projectId }: CreateColumnFormProps) => {
  return (
    <Form className="flex flex-col" method="post" action="/api/columns/create">
      <Label className="mb-4" htmlFor="column-name">
        Name
      </Label>
      <Input
        name="column-name"
        id="column-name"
        type="text"
        min="2"
        max="255"
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
    </Form>
  );
};
