'use client';

import { KeyboardEvent } from 'react';
import { Textarea } from '@/components/ui/textarea';
import { Button } from '@/components/ui/button';
import { IconSend } from '@tabler/icons-react';
import { sendMessage } from './actions';

interface ChatInputProps {
  disabled?: boolean;
  placeholder?: string;
}

export function ChatInput({
  disabled = false,
  placeholder = 'Type a message...',
}: ChatInputProps) {
  const handleKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      e.currentTarget.form?.requestSubmit();
    }
  };

  return (
    <div className="border-t bg-background p-4">
      <form action={sendMessage} className="flex gap-2 items-end">
        <Textarea
          name="message"
          placeholder={placeholder}
          disabled={disabled}
          onKeyDown={handleKeyDown}
          className="min-h-[60px] max-h-[120px] resize-none"
          required
        />
        <Button
          type="submit"
          disabled={disabled}
          size="icon"
          className="h-[35px] w-[35px] flex-shrink-0"
        >
          <IconSend size={20} />
        </Button>
      </form>
      <p className="text-xs text-muted-foreground mt-2">
        Press Enter to send, Shift+Enter for new line
      </p>
    </div>
  );
}
