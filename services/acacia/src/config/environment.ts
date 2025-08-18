import z from 'zod';

const schema = z.object({ port: z.string(), database_url: z.string() });

const unsafeEnvironment = {
  port: process.env.PORT,
  database_url: process.env.DATABASE_URL,
};

const environment = schema.parse(unsafeEnvironment);

export default environment;
