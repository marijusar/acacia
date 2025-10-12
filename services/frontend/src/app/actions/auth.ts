'use server';

import { userService } from '@/lib/services/user-service';
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

  redirect('/projects');
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
