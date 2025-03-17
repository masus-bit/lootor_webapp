import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { ConfigModule } from '@nestjs/config';
import { UserRepository } from '../../repositories/User.repository';
import { User } from '../../entities/User';
import { JwtRegisterModule } from './jwt.module';
import { HttpModule } from '@nestjs/axios';
import { SearchController } from '../../controllers/Search.controller';
import { ElasticsearchService } from '../../services/ElasticSearch.service';
import { EntityModel } from '../../entities/EntityModel';
import { CollectionItem } from '../../entities/CollectionItem';
import { Collection } from '../../entities/Collection';
import { Tags } from '../../entities/Tags';
import { TagsRepository } from '../../repositories/Tags.repository';
import { EntityRepository } from '../../repositories/Entity.repository';
import { CollectionRepository } from '../../repositories/Collection.repository';
import { CollectionItemRepository } from '../../repositories/CollectionItem.repository';

@Module({
  controllers: [SearchController],
  providers: [
    ElasticsearchService,
    UserRepository,
    TagsRepository,
    EntityRepository,
    CollectionRepository,
    CollectionItemRepository,
  ],
  imports: [
    TypeOrmModule.forFeature([
      User,
      EntityModel,
      CollectionItem,
      Collection,
      Tags,
    ]),
    ConfigModule.forRoot({
      envFilePath: '.env',
    }),
    JwtRegisterModule,
    HttpModule,
  ],
  exports: [ElasticsearchService],
})
export class ElasticSearchModule {}
