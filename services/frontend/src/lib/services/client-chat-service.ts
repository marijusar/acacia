type SendMessageStreamArguments = {
  conversationId: number;
  content: string;
  onChunk: (chunk: string) => void;
  onError: (error: Error) => void;
  onDone: () => void;
};

class ClientChatService {
  async sendMessageStream({
    content,
    conversationId,
    onDone,
    onError,
    onChunk,
  }: SendMessageStreamArguments): Promise<void> {
    try {
      const response = await fetch('/api/chat/send-message', {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          conversation_id: conversationId,
          content,
        }),
      });

      if (!response.ok) {
        const body = await response.json();
        const message = `Failed to send message. ${body['message'] || 'Unknown error'}`;
        onError(new Error(message));
        return;
      }

      if (!response.body) {
        onError(new Error('Response body is null'));
        return;
      }

      const reader = response.body.getReader();
      const decoder = new TextDecoder();

      while (true) {
        const { done, value } = await reader.read();

        if (done) break;

        const chunk = decoder.decode(value, { stream: true });
        const lines = chunk.split('\n');

        for (const line of lines) {
          if (line.startsWith('event: ')) {
            const eventType = line.substring(7).trim();

            if (eventType === 'error') {
              // Next line should have the error data
              continue;
            } else if (eventType === 'done') {
              onDone();
              return;
            }
          } else if (line.startsWith('data: ')) {
            const data = line.substring(6);
            if (data) {
              onChunk(data);
            }
          }
        }
      }
    } catch (error) {
      onError(
        error instanceof Error ? error : new Error('Unknown streaming error')
      );
    }
  }
}

export const clientChatService = new ClientChatService();
