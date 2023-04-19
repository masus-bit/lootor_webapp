import { Module } from '@nestjs/common';
import { ExternalService } from '../../services/External.service';
import { HttpModule } from '@nestjs/axios';
import { ExternalController } from '../../controllers/External.controller';
import { ConfigModule } from '@nestjs/config';
import { JwtModule } from '@nestjs/jwt';
import crypto from 'crypto';
import { Algorithm } from 'jsonwebtoken';
import { UserRepository } from '../../repositories/User.repository';
import { TypeOrmModule } from '@nestjs/typeorm';
import { User } from '../../entities/User';

@Module({
  controllers: [ExternalController],
  providers: [ExternalService, UserRepository],
  imports: [
    TypeOrmModule.forFeature([User]),
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
    HttpModule,
  ],
  exports: [ExternalService],
})
export class ExternalModule {}
