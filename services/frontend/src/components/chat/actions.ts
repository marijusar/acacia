'use server';

import {
  createConversationPayload,
  CreateConversationResponse,
} from '@/lib/schemas/conversation';
import { conversationService } from '@/lib/services/conversation-service';
import { revalidatePath } from 'next/cache';

type CreateConversationResult =
  | { error: string; data: null }
  | { error: null; data: CreateConversationResponse };

export async function createConversation(
  _prevState: unknown,
  formData: FormData
): Promise<CreateConversationResult> {
  const entries = Object.fromEntries(formData);
  const requestBody = createConversationPayload.safeParse({
    ...entries,
    initial_message: entries['message'],
  });

  if (requestBody.error) {
    return {
      data: null,
      error: 'Failed to validate input field.',
    };
  }

  const conversation = await conversationService.createConversation(
    requestBody.data
  );

  revalidatePath(`/projects/${requestBody.data.project_id}/board`);

  return { data: conversation, error: null };
}
