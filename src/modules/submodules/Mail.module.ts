import { Module } from '@nestjs/common';
import { MailerModule } from '@nestjs-modules/mailer';

@Module({
  imports: [
    MailerModule.forRoot({
      transport: {
        host: 'smtp.yandex.ru',
        port: 465, // SSL
        secure: true,
        auth: {
          user: 'noreply@lootor.me', // ваша почта
          pass: 'cytuhekbl13', // пароль или app password
        },
      },
      defaults: {
        from: '"Lootor" <noreply@lootor.me>',
      },
    }),
  ],
})
export class MailModule {}
