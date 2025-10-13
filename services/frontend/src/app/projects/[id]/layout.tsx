'use server'

import { projectService } from '@/lib/services/project-service';
import { teamService } from '@/lib/services/team-service';
import { AppSidebar } from '@/components/sidebar/sidebar';
import { Input } from '@/components/ui/input';
import { userService } from '@/lib/services/user-service';
import { Heading1 } from 'lucide-react';
import { redirect } from 'next/navigation';

interface ProjectLayoutProps {
  children: React.ReactNode;
  params: Promise<{ id: string }>;
}

export default async function ProjectLayout({
  children,
  params,
}: ProjectLayoutProps) {
  const { id } = await params;

  const [teams, projectDetails, projects] = await Promise.all([
    teamService.getUserTeams(),
    projectService.getProjectDetails(id),
    projectService.getProjects(),
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
      />
      <div className="flex flex-col flex-1 overflow-hidden">
        <div className="h-16 w-full bg-secondary flex items-center justify-center sticky">
          <Input placeholder="Search.." className="max-w-100" />
        </div>
        <div className="w-24"> </div>
        <div className="flex pt-4 pb-4 pr-8 pl-8 flex-1 flex-col bg-background overflow-auto">
          {children}
        </div>
      </div>
    </div>
  );
}
