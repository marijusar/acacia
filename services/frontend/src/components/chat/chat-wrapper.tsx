'use client';

import { useSearchParams, useRouter } from 'next/navigation';
import { SlidePanel } from '@/components/ui/slide-panel';
import { Chat } from './chat';
import type { ConversationWithMessagesResponse } from '@/lib/schemas/conversation';

interface ChatWrapperProps {
  conversation?: ConversationWithMessagesResponse;
  projectId: number;
}

export function ChatWrapper({ conversation, projectId }: ChatWrapperProps) {
  const searchParams = useSearchParams();
  const router = useRouter();
  const isChatOpen = searchParams.has('chat');

  const handleClose = () => {
    const params = new URLSearchParams(searchParams);
    params.delete('chat');
    router.push(`?${params.toString()}`);
  };

  return (
    <SlidePanel isOpen={isChatOpen} onClose={handleClose} width="400px">
      <Chat
        onClose={handleClose}
        conversation={conversation}
        projectId={projectId}
      />
    </SlidePanel>
  );
}
