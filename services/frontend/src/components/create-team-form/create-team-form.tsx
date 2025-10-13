'use client';

import { useActionState, useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { createTeamAction, type CreateTeamFormState } from '@/app/actions/teams';

const initialState: CreateTeamFormState = {
  error: null,
};

export function CreateTeamForm() {
  const [state, formAction] = useActionState(createTeamAction, initialState);
  const [name, setName] = useState('');

  return (
    <form action={formAction} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="name">Team Name</Label>
        <Input
          id="name"
          name="name"
          type="text"
          placeholder="Enter team name"
          required
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
      </div>
      {state?.error && (
        <div className="text-sm text-destructive">{state.error}</div>
      )}
      <Button type="submit" className="w-full">
        Create
      </Button>
      <div className="text-center text-sm text-muted-foreground">
        Or get invited to one
      </div>
    </form>
  );
}
