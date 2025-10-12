'use client';

import { useActionState, useState } from 'react';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { registerAction, RegisterFormState } from '@/app/actions/auth';

const initialState = {
  error: null,
} satisfies RegisterFormState;

export function RegisterForm() {
  const [state, formAction] = useActionState(registerAction, initialState);
  const [email, setEmail] = useState('');
  const [name, setName] = useState('');
  const [password, setPassword] = useState('');

  return (
    <form action={formAction} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="name">Name</Label>
        <Input
          id="name"
          name="name"
          type="text"
          placeholder="Enter your name"
          required
          autoComplete="name"
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
      </div>
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
          autoComplete="new-password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
      </div>
      {state?.error && (
        <div className="text-sm text-destructive">{state.error}</div>
      )}
      <Button type="submit" className="w-full">
        Register
      </Button>
      <div className="text-center text-sm">
        <span className="text-muted-foreground">
          Already have an account?{' '}
        </span>
        <Link href="/login" className="text-primary hover:underline font-medium">
          Click here to login
        </Link>
      </div>
    </form>
  );
}
