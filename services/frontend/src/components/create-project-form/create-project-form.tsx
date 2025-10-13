'use client';

import { useActionState, useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  createProjectAction,
  type CreateProjectFormState,
} from '@/app/actions/projects';
import { type TeamResponse } from '@/lib/schemas/teams';

const initialState: CreateProjectFormState = {
  error: null,
};

interface CreateProjectFormProps {
  teams: TeamResponse[];
}

export function CreateProjectForm({ teams }: CreateProjectFormProps) {
  const [state, formAction] = useActionState(
    createProjectAction,
    initialState
  );
  const [name, setName] = useState('');
  const [teamId, setTeamId] = useState<string>('');

  return (
    <form action={formAction} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="name">Project Name</Label>
        <Input
          id="name"
          name="name"
          type="text"
          placeholder="Enter project name"
          required
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
      </div>
      <div className="space-y-2">
        <Label htmlFor="team_id">Team</Label>
        <Select name="team_id" value={teamId} onValueChange={setTeamId} required>
          <SelectTrigger id="team_id" className="w-full">
            <SelectValue placeholder="Select a team" />
          </SelectTrigger>
          <SelectContent>
            {teams.map((team) => (
              <SelectItem key={team.id} value={team.id.toString()}>
                {team.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
      {state?.error && (
        <div className="text-sm text-destructive">{state.error}</div>
      )}
      <Button type="submit" className="w-full">
        Create
      </Button>
    </form>
  );
}
