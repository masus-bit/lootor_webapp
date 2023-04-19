import { Module } from '@nestjs/common';
import { JwtModule } from '@nestjs/jwt';

import { Algorithm } from 'jsonwebtoken';
import { ConfigModule, ConfigService } from '@nestjs/config';
import * as crypto from 'crypto';

@Module({
  controllers: [],
  providers: [ConfigService, JwtModule],
  imports: [
    ConfigModule.forRoot({
      envFilePath: '.env',
    }),
    JwtModule.register({
      secret: process.env.JWT_SECRET || crypto.randomBytes(1024).toString(),
      signOptions: {
        algorithm: (process.env.JWT_ALGORITHM as Algorithm) || 'HS512',
        expiresIn: process.env.JWT_EXPIRESS || 86400000,
      },
      verifyOptions: {
        algorithms: [(process.env.JWT_ALGORITHM as Algorithm) || 'HS512'],
      },
    }),
  ],
  exports: [JwtModule],
})
export class JwtRegisterModule {}
