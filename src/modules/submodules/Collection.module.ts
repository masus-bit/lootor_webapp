import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { ConfigModule } from '@nestjs/config';
import { CollectionService } from '../../services/Collection.service';
import { CollectionController } from '../../controllers/Collection.controller';
import { Collection } from '../../entities/Collection';
import { CollectionRepository } from '../../repositories/Collection.repository';
import { UserRepository } from '../../repositories/User.repository';
import { User } from '../../entities/User';
import { EventRepository } from '../../repositories/Event.repository';
import { Event } from '../../entities/Event';
import { EventService } from '../../services/Event.service';
import { JwtRegisterModule } from './jwt.module';
import { HttpModule } from '@nestjs/axios';
import { TagsRepository } from '../../repositories/Tags.repository';
import { TagsService } from '../../services/Tags.service';
import { Tags } from '../../entities/Tags';
import { CollectionItemRepository } from '../../repositories/CollectionItem.repository';
import { CollectionItemService } from '../../services/CollectionItem.service';
import { CollectionItem } from '../../entities/CollectionItem';
import { ElasticsearchService } from '../../services/ElasticSearch.service';
import { EntityRepository } from '../../repositories/Entity.repository';
import { EntityModel } from '../../entities/EntityModel';

@Module({
  controllers: [CollectionController],
  providers: [
    CollectionService,
    CollectionRepository,
    UserRepository,
    EventRepository,
    EventService,
    TagsRepository,
    TagsService,
    CollectionItemRepository,
    CollectionItemService,
    ElasticsearchService,
    EntityRepository,
  ],
  imports: [
    TypeOrmModule.forFeature([
      Collection,
      User,
      Event,
      Tags,
      CollectionItem,
      EntityModel,
    ]),
    ConfigModule.forRoot({
      envFilePath: '.env',
    }),
    JwtRegisterModule,
    HttpModule,
  ],
  exports: [CollectionService],
})
export class CollectionModule {}
