import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { CreateProjectForm } from '@/components/create-project-form/create-project-form';
import { teamService } from '@/lib/services/team-service';

export default async function ProjectsPage() {
  const teams = await teamService.getUserTeams();

  return (
    <div className="flex min-h-screen items-center justify-center">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle>Create project</CardTitle>
          <CardDescription>
            Create a new project to start tracking your work
          </CardDescription>
        </CardHeader>
        <CardContent>
          <CreateProjectForm teams={teams} />
        </CardContent>
      </Card>
    </div>
  );
}
