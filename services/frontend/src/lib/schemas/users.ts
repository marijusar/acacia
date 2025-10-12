import z from 'zod';

export const userResponse = z.object({
  id: z.number(),
  email: z.string().email(),
  name: z.string(),
  created_at: z.string(),
});

export type UserResponse = z.infer<typeof userResponse>;

export const authStatusResponse = z.object({
  authenticated: z.boolean(),
  user: userResponse,
});

export type AuthStatusResponse = z.infer<typeof authStatusResponse>;

export const loginInput = z.object({
  email: z.string().email('Invalid email address'),
  password: z.string().min(1, 'Password is required'),
});

export type LoginInput = z.infer<typeof loginInput>;

export const loginResponse = z.object({
  user: userResponse,
});

export type LoginResponse = z.infer<typeof loginResponse>;

export const registerInput = z.object({
  email: z.string().email('Invalid email address'),
  name: z.string().min(1, 'Name is required'),
  password: z.string().min(6, 'Password must be at least 6 characters'),
});

export type RegisterInput = z.infer<typeof registerInput>;

export const registerResponse = z.object({
  id: z.number(),
  email: z.string(),
  name: z.string(),
  created_at: z.string(),
});

export type RegisterResponse = z.infer<typeof registerResponse>;
