import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { JwtModule } from '@nestjs/jwt';
import { TypeOrmModule } from '@nestjs/typeorm';
import * as crypto from 'crypto';
import { Algorithm } from 'jsonwebtoken';
import { AuthController } from '../controllers/Auth.controller';
import { AuthService } from '../services/Auth.service';
import { UserRepository } from '../repositories/User.repository';
import { User } from '../entities/User';
import { ElasticsearchService } from '../services/ElasticSearch.service';
import { PassportModule } from '@nestjs/passport';

@Module({
  controllers: [AuthController],
  providers: [AuthService, UserRepository, ElasticsearchService],
  imports: [
    ConfigModule.forRoot({
      envFilePath: '.env',
    }),
    PassportModule.register({ defaultStrategy: 'jwt' }),
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
    TypeOrmModule.forFeature([User, UserRepository]),
  ],
  exports: [AuthService, JwtModule],
})
export class AuthModule {}
