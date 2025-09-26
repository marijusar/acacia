import { ChevronsUpDown, AudioWaveform, Plus } from 'lucide-react';

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuTrigger,
} from '../../components/ui/dropdown-menu';
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from '../../components/ui/sidebar';
import type { Project } from '~/schemas/projects';
import { NavLink, useParams, useSearchParams } from 'react-router';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '../../components/ui/dialog';
import { ProjectForm } from '../project-form/project-form';

type ProjectSwitcherProps = {
  projects: Project[];
};

export function ProjectSwitcher({ projects }: ProjectSwitcherProps) {
  const [params, setParams] = useSearchParams();
  const { isMobile } = useSidebar();
  const { id } = useParams();
  const activeProject = projects.find((p) => p.id.toString() == id);
  const createProjectDialog = params.get('create_project_dialog') === 'true';

  if (!activeProject) {
    return null;
  }

  const onClickDialogTrigger = () => {
    console.log('here');
    setParams({ create_project_dialog: 'true' });
  };

  const onDialogOpenChange = () => {
    setParams({ create_project_dialog: `${!createProjectDialog}` });
  };

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <Dialog onOpenChange={onDialogOpenChange} open={!!createProjectDialog}>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <SidebarMenuButton
                size="lg"
                className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground cursor-pointer"
              >
                <div className="bg-sidebar-primary text-sidebar-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg">
                  <Plus className="size-4" />
                </div>
                <div className="grid flex-1 text-left text-sm leading-tight">
                  <span className="truncate font-medium">
                    {activeProject.name}
                  </span>
                  <span className="truncate text-xs">{activeProject.name}</span>
                </div>
                <ChevronsUpDown className="ml-auto" />
              </SidebarMenuButton>
            </DropdownMenuTrigger>
            <DropdownMenuContent
              className="w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"
              align="start"
              side={isMobile ? 'bottom' : 'right'}
              sideOffset={4}
            >
              <DropdownMenuLabel className="text-muted-foreground text-xs">
                Projects
              </DropdownMenuLabel>
              {projects.map((project, index) => (
                <NavLink key={project.id} to={`/projects/board/${project.id}`}>
                  <DropdownMenuItem className="gap-2 p-2 cursor-pointer">
                    <div className="flex size-6 items-center justify-center rounded-md border">
                      <AudioWaveform className="size-3.5 shrink-0" />
                    </div>
                    {project.name}
                    <DropdownMenuShortcut>⌘{index + 1}</DropdownMenuShortcut>
                  </DropdownMenuItem>
                </NavLink>
              ))}
              <DropdownMenuSeparator />
              <DropdownMenuItem
                className="gap-2 p-2"
                onClick={onClickDialogTrigger}
              >
                <div className="flex size-6 items-center justify-center rounded-md border bg-transparent">
                  <Plus className="size-4" />
                </div>
                <DialogTrigger>Create project</DialogTrigger>
              </DropdownMenuItem>
            </DropdownMenuContent>

            <DialogContent aria-describedby={undefined}>
              <DialogHeader>
                <DialogTitle>Project details</DialogTitle>
              </DialogHeader>

              <ProjectForm />
            </DialogContent>
          </DropdownMenu>
        </Dialog>
      </SidebarMenuItem>
    </SidebarMenu>
  );
}
