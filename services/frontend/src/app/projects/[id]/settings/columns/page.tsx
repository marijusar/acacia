import { SettingsCard } from '@/components/settings-card/settings-card';
import { Heading1 } from '@/components/ui/headings';
import { Separator } from '@radix-ui/react-dropdown-menu';
import { DeleteIcon } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { deleteProjectColumnAction } from '@/app/actions/projects';
import { CreateColumnDialog } from '@/components/create-column-dialog/create-column-dialog';
import { projectService } from '@/lib/services/project-service';

interface ColumnSettingsPageProps {
  params: Promise<{ id: string }>;
}

export default async function ColumnSettingsPage({
  params,
}: ColumnSettingsPageProps) {
  const id = (await params).id;
  const columns = await projectService.getProjectColumns(id);

  return (
    <SettingsCard>
      <Heading1>Column Settings</Heading1>
      {columns.map(({ id, name, position_index }) => {
        return (
          <div key={id}>
            <div className="flex">
              <p className="mr-4"> {position_index + 1}. </p>
              <p>{name}</p>

              <form action={deleteProjectColumnAction} className="ml-auto">
                <Input type="hidden" name="project-column-id" value={id} />
                <button
                  type="submit"
                  className="bg-transparent border-none p-0"
                >
                  <DeleteIcon className="stroke-destructive cursor-pointer" />
                </button>
              </form>
            </div>
            <Separator className="w-full bg-secondary-foreground h-0.25 mt-4" />
          </div>
        );
      })}
      <CreateColumnDialog projectId={Number(id)} />
    </SettingsCard>
  );
}
