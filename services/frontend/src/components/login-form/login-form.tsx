'use client';

import { useActionState, useState } from 'react';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { loginAction } from '@/app/actions/auth';

export type LoginFormState = {
  error: string | null;
};

const initialState = {
  error: null,
} satisfies LoginFormState;

export function LoginForm() {
  const [state, formAction] = useActionState(loginAction, initialState);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  return (
    <form action={formAction} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="email">Email</Label>
        <Input
          id="email"
          name="email"
          type="email"
          placeholder="Enter your email"
          required
          autoComplete="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
      </div>
      <div className="space-y-2">
        <Label htmlFor="password">Password</Label>
        <Input
          id="password"
          name="password"
          type="password"
          placeholder="Enter your password"
          required
          autoComplete="current-password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
      </div>
      {state?.error && (
        <div className="text-sm text-destructive">{state.error}</div>
      )}
      <Button type="submit" className="w-full">
        Login
      </Button>
      <div className="text-center text-sm">
        <span className="text-muted-foreground">Don't have an account? </span>
        <Link
          href="/register"
          className="text-primary hover:underline font-medium"
        >
          Click here to register
        </Link>
      </div>
    </form>
  );
}
