import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import crypto from 'crypto';
import { ConfigModule } from '@nestjs/config';
import { JwtModule } from '@nestjs/jwt';
import { Algorithm } from 'jsonwebtoken';
import { DraftLot } from '../../entities/DraftLot';

@Module({
  controllers: [],
  providers: [],
  imports: [
    TypeOrmModule.forFeature([DraftLot]),
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
  exports: [],
})
export class DraftLotModule {}
