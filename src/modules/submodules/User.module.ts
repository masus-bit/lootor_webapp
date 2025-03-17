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
import { CollectionRepository } from '../../repositories/Collection.repository';
import { CollectionService } from '../../services/Collection.service';
import { CollectionItemService } from '../../services/CollectionItem.service';
import { CollectionItemRepository } from '../../repositories/CollectionItem.repository';
import { Collection } from '../../entities/Collection';
import { CollectionItem } from '../../entities/CollectionItem';
import { ElasticsearchService } from '../../services/ElasticSearch.service';
import { EntityRepository } from '../../repositories/Entity.repository';
import { EntityModel } from '../../entities/EntityModel';

@Module({
  controllers: [UserController],
  providers: [
    UserService,
    UserRepository,
    EventRepository,
    CollectionRepository,
    CollectionService,
    CollectionItemService,
    CollectionItemRepository,
    ElasticsearchService,
    EntityRepository,
  ],
  imports: [
    TypeOrmModule.forFeature([
      User,
      Event,
      Collection,
      CollectionItem,
      EntityModel,
    ]),
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
