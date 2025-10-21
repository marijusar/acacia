'use server';

import { projectService } from '@/lib/services/project-service';
import { teamService } from '@/lib/services/team-service';
import { userService } from '@/lib/services/user-service';
import { conversationService } from '@/lib/services/conversation-service';
import { AppSidebar } from '@/components/sidebar/sidebar';
import { Input } from '@/components/ui/input';
import { Heading1 } from 'lucide-react';
import { redirect } from 'next/navigation';
import { ChatWrapper } from '@/components/chat';

interface ProjectLayoutProps {
  children: React.ReactNode;
  params: Promise<{ id: string }>;
}

export default async function ProjectLayout({
  children,
  params,
}: ProjectLayoutProps) {
  const { id } = await params;

  const [teams, projectDetails, projects, authStatus, latestConversation] =
    await Promise.all([
      teamService.getUserTeams(),
      projectService.getProjectDetails(id),
      projectService.getProjects(),
      userService.getAuthStatus(),
      conversationService.getLatestConversation(),
    ]);

  // Check if user has teams
  if (!teams || teams.length === 0) {
    redirect('/teams');
  }

  // Check if user has projects
  if (!projects || projects.length === 0) {
    redirect('/');
  }

  if (!projects || !projectDetails) {
    return <Heading1>Not implemented</Heading1>;
  }

  return (
    <div className="flex h-full overflow-y-hidden">
      <AppSidebar
        projects={projects}
        projectName={projectDetails.name}
        projectId={projectDetails.id}
        teams={teams}
        user={authStatus.user}
      />
      <div className="flex flex-col flex-1 overflow-hidden">
        <div className="h-16 w-full bg-secondary flex items-center justify-center sticky"></div>
        <div className="w-24"> </div>
        <div className="flex flex-1">
          <div className="flex pt-4 pb-4 pr-8 pl-8 flex-1 flex-col bg-background overflow-auto">
            {children}
          </div>
          <ChatWrapper
            conversation={latestConversation}
            projectId={projectDetails.id}
          />
        </div>
      </div>
    </div>
  );
}
