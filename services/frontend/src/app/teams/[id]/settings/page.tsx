import { redirect } from 'next/navigation';

interface TeamsSettingsPageProps {
  params: Promise<{ id: string }>;
}

export default async function TeamsSettingsPage({
  params,
}: TeamsSettingsPageProps) {
  const { id } = await params;
  redirect(`/teams/${id}/settings/api-keys`);
}
