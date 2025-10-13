'use server';

import { userService } from '@/lib/services/user-service';
import { teamService } from '@/lib/services/team-service';
import { projectService } from '@/lib/services/project-service';
import { loginInput, registerInput } from '@/lib/schemas/users';
import { redirect } from 'next/navigation';
import { LoginFormState } from '@/components/login-form/login-form';

export async function loginAction(
  _: LoginFormState,
  formData: FormData
): Promise<LoginFormState> {
  try {
    const credentials = {
      email: formData.get('email'),
      password: formData.get('password'),
    };

    const validatedCredentials = loginInput.parse(credentials);

    await userService.login(validatedCredentials);
  } catch (error) {
    if (error instanceof Error) {
      return {
        error: error.message,
      };
    }
    return {
      error: 'An unexpected error has occurred',
    };
  }

  // Fetch user's teams and projects in parallel
  const [teams, projects] = await Promise.all([
    teamService.getUserTeams(),
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

  // Redirect to first project's board
  redirect(`/projects/${projects[0].id}/board`);
}

export type RegisterFormState = {
  error: string | null;
};

export async function registerAction(
  _: RegisterFormState,
  formData: FormData
): Promise<RegisterFormState> {
  try {
    const data = {
      email: formData.get('email'),
      name: formData.get('name'),
      password: formData.get('password'),
    };

    const validatedData = registerInput.parse(data);

    await userService.register(validatedData);
  } catch (error) {
    if (error instanceof Error) {
      return {
        error: error.message,
      };
    }
    return {
      error: 'An unexpected error has occurred',
    };
  }

  redirect('/login');
}
