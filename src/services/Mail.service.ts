import { Injectable } from '@nestjs/common';
import { MailerService } from '@nestjs-modules/mailer';

@Injectable()
export class MailService {
  constructor(private readonly mailerService: MailerService) {}

  async sendConfirmationEmail(email: string, token: string) {
    await this.mailerService.sendMail({
      to: email,
      subject: 'Подтверждение регистрации на Lootor',
      text: `Подтвердите ваш аккаунт, перейдя по ссылке: https://dev.lootor.me/auth?confirmationToken=${token}`,
    });
  }
}
