import { Cog, Inbox, MessageCircle } from 'lucide-react';
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarProvider,
} from '../ui/sidebar';
import { NavUser } from '../nav-user/nav-user';
import { ProjectSwitcher } from '../project-switcher/project-switcher';
import type { Project } from '@/lib/schemas/projects';
import type { TeamResponse } from '@/lib/schemas/teams';
import type { UserResponse } from '@/lib/schemas/users';
import { CollapsibleItem } from './collapsible-item';
import { DefaultItem } from './default-item';

const items = [
  {
    title: 'Board',
    url: (id: number) => `/projects/${id}/board`,
    icon: Inbox,
  },

  {
    title: 'Chat',
    url: (id: number) => `/projects/${id}/board?chat`,
    icon: MessageCircle,
  },
  {
    title: 'Settings',
    url: (id: number) => `/projects/${id}/settings`,
    icon: Cog,
    items: [
      {
        title: 'Columns',
        url: (id: number) => `/projects/${id}/settings/columns`,
      },
    ],
  },
];

type AppSidebarProps = {
  projectName: string;
  projectId: number;
  projects: Project[];
  teams: TeamResponse[];
  user: UserResponse;
};

export function AppSidebar({
  projectName,
  projectId,
  projects,
  teams,
  user,
}: AppSidebarProps) {
  return (
    <SidebarProvider className="max-w-(--sidebar-width)">
      <Sidebar>
        <SidebarContent>
          <SidebarHeader>
            <ProjectSwitcher projects={projects} />
          </SidebarHeader>
          <SidebarGroup>
            <SidebarGroupLabel>{projectName}</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {items.map((item, i) =>
                  item.items && item.items.length > 0 ? (
                    <CollapsibleItem {...item} projectId={projectId} key={i} />
                  ) : (
                    <DefaultItem
                      {...item}
                      icon={<item.icon />}
                      projectId={projectId}
                      key={i}
                    />
                  )
                )}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>
        <SidebarFooter>
          <NavUser
            user={{
              ...user,
              avatar: 'https://github.com/shadcn.png',
            }}
            teams={teams}
          />
        </SidebarFooter>
      </Sidebar>
    </SidebarProvider>
  );
}
