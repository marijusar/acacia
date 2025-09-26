import { SettingsCard } from '~/components/settings-card/settings-card';
import { Heading1 } from '~/components/ui/headings';
import type { Route } from './+types/column-settings';
import { Separator } from '@radix-ui/react-dropdown-menu';
import { DeleteIcon, PlusIcon } from 'lucide-react';
import { useSearchParams, useFetcher } from 'react-router';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '~/components/ui/dialog';
import { CreateColumnForm } from '~/components/create-column-form/create-column-form';
import { useDashboardContext } from '~/layouts/dashboard';
import { Input } from '~/components/ui/input';

const ColumnSettings = ({}: Route.ComponentProps) => {
  const {
    projectDetails: { columns, id },
  } = useDashboardContext();
  const [params, setParams] = useSearchParams();
  const fetcher = useFetcher();

  const onToggleCreateColumnModal = () => {
    setParams((params) => {
      params.set(
        'create_column',
        String(!(params.get('create_column') === 'true'))
      );

      return params;
    });
  };

  return (
    <SettingsCard>
      <Heading1>Column Settings</Heading1>
      {columns.map(({ id, name, position_index }) => {
        return (
          <div key={id}>
            <div className="flex">
              <p className="mr-4"> {position_index + 1}. </p>
              <p>{name}</p>

              <fetcher.Form
                method="post"
                action="/api/columns/remove"
                className="ml-auto"
              >
                <Input type="hidden" name="project-column-id" value={id} />
                <button
                  type="submit"
                  className="bg-transparent border-none p-0"
                >
                  <DeleteIcon className="stroke-destructive cursor-pointer" />
                </button>
              </fetcher.Form>
            </div>
            <Separator className="w-full bg-secondary-foreground h-0.25 mt-4" />
          </div>
        );
      })}
      <PlusIcon
        onClick={onToggleCreateColumnModal}
        className="bg-sidebar-primary rounded-sm ml-auto cursor-pointer"
      />

      <Dialog
        onOpenChange={onToggleCreateColumnModal}
        open={params.get('create_column') === 'true'}
      >
        <DialogContent aria-describedby={undefined}>
          <DialogHeader>
            <DialogTitle>Create column</DialogTitle>
          </DialogHeader>

          <CreateColumnForm projectId={id} />
        </DialogContent>
      </Dialog>
    </SettingsCard>
  );
};

export default ColumnSettings;
