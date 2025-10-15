'use client';

import { IconBrandOpenai, IconKey } from '@tabler/icons-react';
import type { TeamLLMAPIKeysListResponse } from '@/lib/schemas/team-llm-api-keys';
import { APIKeyEntry } from './api-key-entry';

interface APIKeysManagerProps {
  teamId: number;
  initialApiKeys: TeamLLMAPIKeysListResponse;
}

const PROVIDERS = [
  {
    id: 'anthropic',
    name: 'Anthropic',
    icon: IconKey,
  },
  {
    id: 'openai',
    name: 'OpenAI',
    icon: IconBrandOpenai,
  },
] as const;

export function APIKeysManager({
  teamId,
  initialApiKeys,
}: APIKeysManagerProps) {
  const getExistingKey = (providerId: string) => {
    return initialApiKeys.find((key) => key.provider === providerId);
  };

  return (
    <div className="space-y-6">
      {PROVIDERS.map((provider) => (
        <APIKeyEntry
          key={provider.id}
          providerId={provider.id}
          providerName={provider.name}
          Icon={provider.icon}
          teamId={teamId}
          existingKey={getExistingKey(provider.id)}
        />
      ))}
    </div>
  );
}
