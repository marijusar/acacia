import { cn } from '@/lib/utils';
import type { ChatMessage } from './chat';
import { useEffect, useState } from 'react';

interface ChatMessageProps {
  message: ChatMessage;
}

export function ChatMessage({ message }: ChatMessageProps) {
  const isUser = message.role === 'user';
  const [timestamp, setTimestamp] = useState<string | null>(null);

  useEffect(() => {
    // Since user's and server's timestamps are different, render timestamp only on the client side to avoid hydration errors..
    if (!timestamp) {
      setTimestamp(
        message.timestamp.toLocaleTimeString([], {
          hour: '2-digit',
          minute: '2-digit',
        })
      );
    }
  }, [message.timestamp, timestamp]);

  return (
    <div
      className={cn(
        'flex w-full mb-4 min-h-[56px]',
        isUser ? 'justify-end' : 'justify-start'
      )}
    >
      <div
        className={cn(
          'max-w-[80%] min-w-[100px] rounded-lg px-4 py-2 flex flex-col',
          isUser
            ? 'bg-primary text-primary-foreground'
            : 'bg-muted text-foreground'
        )}
      >
        <p className="text-sm whitespace-pre-wrap break-words">
          {message.content}
        </p>
        <span className="text-xs opacity-70 mt-1 block mt-auto">
          {timestamp}
        </span>
      </div>
    </div>
  );
}
