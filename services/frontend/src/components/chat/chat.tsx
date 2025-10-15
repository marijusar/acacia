'use client';

import { useState } from 'react';
import { ChatMessages } from './chat-messages';
import { ChatInput } from './chat-input';
import { Button } from '@/components/ui/button';
import { IconX } from '@tabler/icons-react';
import type { Message } from './types';

interface ChatProps {
  onClose?: () => void;
}

export function Chat({ onClose }: ChatProps) {
  const [messages] = useState<Message[]>(
    Array(100)
      .fill(null)
      .map((c, i) => ({
        role: 'user',
        content: 'lololol',
        timestamp: new Date(),
        id: i.toString(),
      }))
  );

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="border-b bg-background p-4 flex items-center justify-between">
        <h2 className="font-semibold">Project Assistant</h2>
        {onClose && (
          <Button variant="ghost" size="icon" onClick={onClose}>
            <IconX size={20} />
          </Button>
        )}
      </div>

      {/* Messages */}
      <ChatMessages messages={messages} />

      {/* Input */}
      <ChatInput placeholder="Ask about this project..." />
    </div>
  );
}
