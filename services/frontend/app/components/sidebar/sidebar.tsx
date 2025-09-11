import { Inbox } from 'lucide-react';
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
} from '../ui/sidebar';
import { Avatar, AvatarImage } from '../ui/avatar';
import { NavUser } from '../nav-user/nav-user';
const items = [
  {
    title: 'Board',
    url: '#',
    icon: Inbox,
  },
];

type AppSidebarProps = {
  name: string;
};

export function AppSidebar({ name }: AppSidebarProps) {
  return (
    <SidebarProvider className="max-w-(--sidebar-width)">
      <Sidebar>
        <SidebarContent>
          <SidebarHeader></SidebarHeader>
          <SidebarGroup>
            <SidebarGroupLabel>{name}</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {items.map((item) => (
                  <SidebarMenuItem key={item.title}>
                    <SidebarMenuButton asChild>
                      <a className="bg-primary" href={item.url}>
                        <item.icon />
                        <span>{item.title}</span>
                      </a>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>
        <SidebarFooter>
          <NavUser
            user={{
              name: 'marijus',
              email: 'marijus@acacia.com',
              avatar: 'https://github.com/shadcn.png',
            }}
          />
        </SidebarFooter>
      </Sidebar>
    </SidebarProvider>
  );
}
