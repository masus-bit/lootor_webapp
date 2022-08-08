import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import crypto from 'crypto';
import { ConfigModule } from '@nestjs/config';
import { JwtModule } from '@nestjs/jwt';
import { Algorithm } from 'jsonwebtoken';
import { CollectionService } from '../../services/Collection.service';
import { CollectionController } from '../../controllers/Collection.controller';
import { Collection } from '../../entities/Collection';
import { CollectionRepository } from '../../repositories/Collection.repository';
import { UserRepository } from '../../repositories/User.repository';
import { User } from '../../entities/User';

@Module({
  controllers: [CollectionController],
  providers: [CollectionService, CollectionRepository, UserRepository],
  imports: [
    TypeOrmModule.forFeature([Collection, User]),
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
  exports: [CollectionService],
})
export class CollectionModule {}
