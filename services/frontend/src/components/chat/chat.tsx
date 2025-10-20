'use client';

import { ChatMessages } from './chat-messages';
import { ChatInput } from './chat-input';
import { Button } from '@/components/ui/button';
import { IconX } from '@tabler/icons-react';
import type { ConversationWithMessagesResponse } from '@/lib/schemas/conversation';
import { useChatMessages } from '@/hooks/use-chat-messages';

interface ChatProps {
  onClose?: () => void;
  conversation?: ConversationWithMessagesResponse;
  projectId: number;
}

export type ChatMessage = {
  role: string;
  content: string;
  timestamp: Date;
  streaming: boolean;
};

export function Chat({ onClose, conversation, projectId }: ChatProps) {
  const { messages, handleSendMessage } = useChatMessages(conversation);

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="border-b bg-background p-4 flex items-center justify-between">
        <h2 className="font-semibold">
          {conversation ? conversation.conversation.title : 'Project assistant'}
        </h2>
        {onClose && (
          <Button variant="ghost" size="icon" onClick={onClose}>
            <IconX size={20} />
          </Button>
        )}
      </div>

      {/* Messages */}
      <ChatMessages messages={messages} />

      {/* Input */}
      <ChatInput
        placeholder="Ask about this project..."
        projectId={projectId}
        onMessageSend={handleSendMessage}
      />
    </div>
  );
}
