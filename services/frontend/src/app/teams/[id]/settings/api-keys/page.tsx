import { teamLLMAPIKeysService } from '@/lib/services/team-llm-api-keys-service';
import { APIKeysManager } from '@/components/team-settings/api-keys-manager';
import { SettingsCard } from '@/components/settings-card/settings-card';
import { Heading1 } from '@/components/ui/headings';

interface APIKeysPageProps {
  params: Promise<{ id: string }>;
}

export default async function APIKeysPage({ params }: APIKeysPageProps) {
  const { id } = await params;
  const teamId = parseInt(id);

  const apiKeys = await teamLLMAPIKeysService.getAPIKeys(teamId);

  return (
    <SettingsCard>
      <Heading1>API Keys</Heading1>
      <p className="text-muted-foreground mb-8">
        Manage your LLM provider API keys. Keys are encrypted and stored
        securely.
      </p>

      <APIKeysManager teamId={teamId} initialApiKeys={apiKeys} />
    </SettingsCard>
  );
}
