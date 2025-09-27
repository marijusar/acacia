import type { ReactNode } from 'react';
import { Card } from '../ui/card';

type SettingsCardProps = {
  children: ReactNode;
};

export const SettingsCard = ({ children }: SettingsCardProps) => {
  return <Card className="flex-1 p-8">{children}</Card>;
};
