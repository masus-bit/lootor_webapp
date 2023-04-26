import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { User } from '../../entities/User';
import { UserController } from '../../controllers/User.controller';
import { UserService } from '../../services/User.service';
import { UserRepository } from '../../repositories/User.repository';
import crypto from 'crypto';
import { ConfigModule } from '@nestjs/config';
import { JwtModule } from '@nestjs/jwt';
import { Algorithm } from 'jsonwebtoken';
import { EventRepository } from '../../repositories/Event.repository';
import { Event } from '../../entities/Event';

@Module({
  controllers: [UserController],
  providers: [UserService, UserRepository, EventRepository],
  imports: [
    TypeOrmModule.forFeature([User, Event]),
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
  exports: [UserService],
})
export class UserModule {}
