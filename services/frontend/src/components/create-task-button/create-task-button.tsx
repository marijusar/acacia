'use client';

import { useRouter } from 'next/navigation';
import { Button } from '../ui/button';
import { usePathname, useSearchParams } from 'next/navigation';

export const CreateTaskButton = () => {
  const router = useRouter();
  const params = useSearchParams();
  const pathname = usePathname();
  const route = `${pathname}?${params.toString()}&open_issue_id=new`;
  return (
    <Button onClick={() => router.push(route)} className="ml-auto">
      Create task
    </Button>
  );
};
