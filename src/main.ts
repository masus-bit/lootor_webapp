import { AppModule } from './App.module';
import { NestFactory } from '@nestjs/core';
import { ValidationPipe } from '@nestjs/common';
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger';

(async () => {
  const PORT = process.env.PORT || 5000;
  // const options = {
  //   cors: true,
  // };

  const app = await NestFactory.create(AppModule);
  const config = new DocumentBuilder()
    .setTitle('auc documentation')
    .setDescription('Дока для ауса')
    .setVersion('0.0.1')
    .build();
  const document = SwaggerModule.createDocument(app, config);
  SwaggerModule.setup('/api/v0/docs', app, document);
  app.useGlobalPipes(new ValidationPipe());
  app.enableCors({
    origin: '*',})
  await (
    await app
  ).listen(PORT, '0.0.0.0', () =>
    console.log(`server started at хуерт ${PORT}`),
  );
})();
