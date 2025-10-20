import { createConversation } from '@/components/chat/actions';
import { ChatMessage } from '@/components/chat/chat';
import { logger } from '@/lib/config/logger';
import { ConversationWithMessagesResponse } from '@/lib/schemas/conversation';
import { clientChatService } from '@/lib/services/client-chat-service';
import { useState } from 'react';

export const useChatMessages = (
  conversation: ConversationWithMessagesResponse | undefined
) => {
  const [messages, setMessages] = useState<ChatMessage[]>(() => {
    if (!conversation) return [];
    return conversation.messages.map((msg) => ({
      role: msg.role,
      content: msg.content,
      timestamp: new Date(msg.created_at),
      streaming: false,
    }));
  });

  const conversationId = conversation?.conversation.id;

  const handleSendMessage = async (form: FormData) => {
    const content = form.get('message')?.toString() || '';

    let cid = conversationId;
    if (!cid) {
      const result = await createConversation({}, form);
      if (result.error !== null) {
        logger.error('[CHAT] Failed to create conversation', {
          error: result.error,
        });
        throw new Error(result.error);
      }
      cid = result.data.id;
    }

    const userMessage: ChatMessage = {
      role: 'user',
      content,
      timestamp: new Date(),
      streaming: false,
    };

    setMessages((prev) => [...prev, userMessage]);

    const assistantMessage: ChatMessage = {
      role: 'assistant',
      content: '',
      timestamp: new Date(),
      streaming: true,
    };

    setMessages((prev) => [...prev, assistantMessage]);

    // Stream the response
    try {
      await clientChatService.sendMessageStream({
        conversationId: cid,
        content,
        onChunk: (chunk) => {
          // Update the streaming message with new content
          setMessages((prev) => {
            const lastIndex = prev.length - 1;
            return [
              ...prev.slice(0, -1),
              { ...prev[lastIndex], content: prev[lastIndex].content + chunk },
            ];
          });
        },
        onError: (error) => {
          logger.error('Streaming error:', error);
          // Mark streaming as complete even on error
          setMessages((prev) => {
            const lastIndex = prev.length - 1;
            const curr = [...prev];
            curr[lastIndex].streaming = false;
            return prev;
          });
        },
        onDone: () => {
          setMessages((prev) => {
            const lastIndex = prev.length - 1;
            const curr = [...prev];
            curr[lastIndex].streaming = false;
            return prev;
          });
        },
      });
    } catch (error) {
      logger.error('Failed to send message:', { error });
      setMessages((prev) => prev.slice(0, -1));
    }
  };

  return {
    messages,
    handleSendMessage,
  };
};
