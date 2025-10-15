'use server';

import { projectService } from '@/lib/services/project-service';
import { teamService } from '@/lib/services/team-service';
import { userService } from '@/lib/services/user-service';
import { AppSidebar } from '@/components/sidebar/sidebar';
import { Input } from '@/components/ui/input';
import { redirect } from 'next/navigation';

interface TeamSettingsLayoutProps {
  children: React.ReactNode;
  params: Promise<{ id: string }>;
}

export default async function TeamSettingsLayout({
  children,
  params,
}: TeamSettingsLayoutProps) {
  const { id } = await params;

  // Verify user has access to this team
  const [teams, projects, authStatus] = await Promise.all([
    teamService.getUserTeams(),
    projectService.getProjects(),
    userService.getAuthStatus(),
  ]);

  const team = teams.find((t) => t.id === parseInt(id));

  if (!team) {
    redirect('/teams');
  }

  // Check if user has projects for the sidebar
  if (!projects || projects.length === 0) {
    redirect('/');
  }

  return (
    <div className="flex h-full overflow-y-hidden">
      <AppSidebar
        projects={projects}
        projectName={team.name}
        projectId={team.id}
        teams={teams}
        user={authStatus.user}
      />
      <div className="flex flex-col flex-1 overflow-hidden">
        <div className="h-16 w-full bg-secondary flex items-center justify-center sticky">
          <Input placeholder="Search.." className="max-w-100" />
        </div>
        <div className="flex flex-1">
          <div className="flex pt-4 pb-4 pr-8 pl-8 flex-1 flex-col bg-background overflow-auto">
            {children}
          </div>
        </div>
      </div>
    </div>
  );
}
