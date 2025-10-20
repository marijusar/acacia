import z from 'zod';

export const createConversationPayload = z.object({
  initial_message: z.string(),
  model: z.string(),
  provider: z.string(),
  project_id: z.string().transform((v) => parseInt(v)),
});

export type CreateConversationPayload = z.infer<
  typeof createConversationPayload
>;

export const createConversationResponse = z.object({
  id: z.number(),
  user_id: z.number(),
  title: z.string(),
  provider: z.string(),
  model: z.string(),
  created_at: z.string(),
  updated_at: z.string(),
});

export type CreateConversationResponse = z.infer<
  typeof createConversationResponse
>;

export const messageResponse = z.object({
  id: z.number(),
  conversation_id: z.number(),
  role: z.enum(['user', 'assistant']),
  content: z.string(),
  sequence_number: z.number(),
  created_at: z.string(),
});

export type MessageResponse = z.infer<typeof messageResponse>;

export const conversationWithMessagesResponse = z.object({
  conversation: createConversationResponse,
  messages: z.array(messageResponse),
});

export type ConversationWithMessagesResponse = z.infer<
  typeof conversationWithMessagesResponse
>;

export const sendMessagePayload = z.object({
  conversation_id: z.number(),
  content: z.string(),
});

export type SendMessagePayload = z.infer<typeof sendMessagePayload>;
