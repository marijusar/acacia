import type { ReactNode } from 'react';
import { SidebarMenuSubButton, SidebarMenuSubItem } from '../ui/sidebar';

type DefaultItemProps = {
  title: string;
  projectId: number;
  icon: ReactNode;
  url: (arg: number) => string;
};

export const DefaultItem = ({
  title,
  projectId,
  url,
  icon,
}: DefaultItemProps) => {
  return (
    <SidebarMenuSubItem key={title}>
      <SidebarMenuSubButton asChild>
        <a href={url(projectId)}>
          {icon}
          {title}
        </a>
      </SidebarMenuSubButton>
    </SidebarMenuSubItem>
  );
};
