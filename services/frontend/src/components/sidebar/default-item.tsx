import type { ReactNode } from 'react';
import { SidebarMenuSubButton, SidebarMenuSubItem } from '../ui/sidebar';
import Link from 'next/link';

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
        <Link href={url(projectId)}>
          {icon}
          {title}
        </Link>
      </SidebarMenuSubButton>
    </SidebarMenuSubItem>
  );
};
