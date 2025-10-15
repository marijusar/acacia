'use client';

import { useActionState } from 'react';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { IconTrash } from '@tabler/icons-react';
import type { Icon as TablerIcon } from '@tabler/icons-react';
import {
  saveAPIKeyAction,
  deleteAPIKeyAction,
} from '@/app/teams/[id]/settings/api-keys/actions';

interface APIKeyEntryProps {
  providerId: string;
  providerName: string;
  Icon: typeof TablerIcon;
  teamId: number;
  existingKey?: {
    id: number;
    provider: string;
  };
}

export function APIKeyEntry({
  providerId,
  providerName,
  Icon,
  teamId,
  existingKey,
}: APIKeyEntryProps) {
  const [saveState, saveFormAction, isSaving] = useActionState(
    saveAPIKeyAction,
    {
      success: true,
    }
  );

  const handleDelete = async () => {
    if (!existingKey) return;

    if (
      !window.confirm(
        `Are you sure you want to delete the ${providerName} API key?`
      )
    ) {
      return;
    }

    await deleteAPIKeyAction(teamId, existingKey.id);
  };

  return (
    <div className="border rounded-lg p-6 bg-card space-y-4">
      <div className="flex items-center gap-3">
        <div className="flex items-center justify-center w-10 h-10 rounded-lg bg-primary/10">
          <Icon className="w-6 h-6" />
        </div>
        <div>
          <h3 className="font-semibold">{providerName}</h3>
          <p className="text-sm text-muted-foreground">
            {existingKey ? 'API key configured' : 'No API key configured'}
          </p>
        </div>
      </div>

      <form action={saveFormAction} className="space-y-4">
        <input type="hidden" name="team_id" value={teamId} />
        <input type="hidden" name="provider" value={providerId} />

        <div className="flex gap-2">
          <Input
            type="password"
            name="api_key"
            placeholder={
              existingKey
                ? '********************************'
                : `Enter ${providerName} API key`
            }
            disabled={isSaving}
            className="flex-1"
            required
          />
          <Button type="submit" disabled={isSaving}>
            {existingKey ? 'Update' : 'Save'}
          </Button>
          {existingKey && (
            <Button
              type="button"
              variant="destructive"
              size="icon"
              onClick={handleDelete}
              disabled={isSaving}
            >
              <IconTrash className="w-4 h-4" />
            </Button>
          )}
        </div>

        {!saveState.success && saveState.error && (
          <p className="text-sm text-destructive">{saveState.error}</p>
        )}
      </form>
    </div>
  );
}
