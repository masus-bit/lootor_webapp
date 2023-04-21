import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import crypto from 'crypto';
import { ConfigModule } from '@nestjs/config';
import { JwtModule } from '@nestjs/jwt';
import { Algorithm } from 'jsonwebtoken';
import { UserRepository } from '../../repositories/User.repository';
import { User } from '../../entities/User';
import { CollectionItemController } from '../../controllers/CollectionItem.controller';
import { CollectionItemService } from '../../services/CollectionItem.service';
import { CollectionItemRepository } from '../../repositories/CollectionItem.repository';
import { CollectionItem } from '../../entities/CollectionItem';

@Module({
  controllers: [CollectionItemController],
  providers: [CollectionItemService, CollectionItemRepository, UserRepository],
  imports: [
    TypeOrmModule.forFeature([CollectionItem, User]),
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
  exports: [CollectionItemService],
})
export class CollectionItemModule {}
