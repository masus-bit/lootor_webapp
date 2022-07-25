import { AppModule } from './App.module';
import { NestFactory } from '@nestjs/core';

(async () => {
  const PORT = process.env.PORT || 5000;
  const options = {
    cors: true,
  };
  const app = await NestFactory.create(AppModule, options);
  await (
    await app
  ).listen(PORT, '0.0.0.0', () =>
    console.log(`server started at port ${PORT}`),
  );
})();
