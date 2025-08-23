import { testDatabaseContainer } from './database-setup.ts';

export const setup = () => {
  testDatabaseContainer.verify();
};

export const teardown = async () => {
  await testDatabaseContainer.container.stop();
};
