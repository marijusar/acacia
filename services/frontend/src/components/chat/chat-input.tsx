'use client';

import { KeyboardEvent } from 'react';
import { Textarea } from '@/components/ui/textarea';
import { Button } from '@/components/ui/button';
import { IconSend } from '@tabler/icons-react';

interface ChatInputProps {
  disabled?: boolean;
  placeholder?: string;
  projectId: number;
  onMessageSend: (formData: FormData) => void;
}

export function ChatInput({
  disabled = false,
  placeholder = 'Type a message...',
  projectId,
  onMessageSend,
}: ChatInputProps) {
  const handleKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      e.currentTarget.form?.requestSubmit();
    }
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    onMessageSend(formData);
    e.currentTarget.reset();
  };

  return (
    <div className="border-t bg-background p-4">
      <form onSubmit={handleSubmit} className="flex gap-2 items-end">
        <Textarea
          name="message"
          placeholder={placeholder}
          disabled={disabled}
          onKeyDown={handleKeyDown}
          className="min-h-[60px] max-h-[120px] resize-none"
          required
        />
        <input type="hidden" name="model" value="gpt-4o-mini" />
        <input type="hidden" name="provider" value="openai" />
        <input type="hidden" name="project_id" value={projectId} />
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
