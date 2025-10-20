import { env } from '../config/env';
import { logger } from '../config/logger';
import {
  CreateConversationPayload,
  createConversationResponse,
  conversationWithMessagesResponse,
  ConversationWithMessagesResponse,
} from '../schemas/conversation';
import { BaseHttpService, BaseServiceArguments } from './base-service';

type ConversationServiceArguments = BaseServiceArguments & {};

class ConversationService extends BaseHttpService {
  constructor(args: ConversationServiceArguments) {
    super(args);
  }

  async createConversation(params: CreateConversationPayload) {
    const response = await fetch(`${this.url}/conversations`, {
      body: JSON.stringify(params),
      method: 'POST',
      headers: {
        Cookie: await this.cookieService.getAuthCookies(),
      },
    });

    if (!response.ok) {
      const body = await response.json();
      const message = `Failed to create conversation. ${body['message']}`;
      this.logger.error(message, { ...response, body });
      throw new Error(message);
    }

    const body = await response.json();

    return createConversationResponse.parse(body);
  }

  async getLatestConversation(): Promise<
    ConversationWithMessagesResponse | undefined
  > {
    const response = await fetch(`${this.url}/conversations/latest`, {
      method: 'GET',
      headers: {
        Cookie: await this.cookieService.getAuthCookies(),
      },
    });

    if (response.status === 404) {
      // No conversations found - this is expected for new users
      return undefined;
    }

    if (!response.ok) {
      const body = await response.json();
      const message = `Failed to get latest conversation. ${body['message']}`;
      this.logger.error(message, { ...response, body });
      throw new Error(message);
    }

    const body = await response.json();

    return conversationWithMessagesResponse.parse(body);
  }

  async sendMessageStream(
    conversationId: number,
    content: string,
    controller: ReadableStreamDefaultController
  ): Promise<void> {
    const encoder = new TextEncoder();

    try {
      const response = await fetch(`${this.url}/conversations/messages`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Cookie: await this.cookieService.getAuthCookies(),
        },
        body: JSON.stringify({
          conversation_id: conversationId,
          content,
        }),
      });

      if (!response.ok) {
        const body = await response.json();
        const message = `Failed to send message. ${body['message']}`;
        this.logger.error(message, { ...response, body });
        controller.enqueue(
          encoder.encode(`event: error\ndata: ${message}\n\n`)
        );
        controller.close();
        return;
      }

      if (!response.body) {
        controller.enqueue(
          encoder.encode(`event: error\ndata: Response body is null\n\n`)
        );
        controller.close();
        return;
      }

      const reader = response.body.getReader();

      while (true) {
        const { done, value } = await reader.read();

        if (done) {
          controller.close();
          break;
        }

        // Forward the chunk directly to the controller
        controller.enqueue(value);
      }
    } catch (error) {
      this.logger.error('Error reading stream', { error });
      const message =
        error instanceof Error ? error.message : 'Unknown streaming error';
      controller.enqueue(encoder.encode(`event: error\ndata: ${message}\n\n`));
      controller.close();
    }
  }
}

export const conversationService = new ConversationService({
  logger,
  url: env.ACACIA_API_URL,
});
