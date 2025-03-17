import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { ConfigModule } from '@nestjs/config';
import { UserRepository } from '../../repositories/User.repository';
import { User } from '../../entities/User';
import { JwtRegisterModule } from './jwt.module';
import { HttpModule } from '@nestjs/axios';
import { Collection } from '../../entities/Collection';
import { TagsService } from '../../services/Tags.service';
import { TagsRepository } from '../../repositories/Tags.repository';
import { TagsController } from '../../controllers/Tags.controller';
import { CollectionRepository } from '../../repositories/Collection.repository';
import { Tags } from '../../entities/Tags';
import { ElasticsearchService } from '../../services/ElasticSearch.service';
import { EntityRepository } from '../../repositories/Entity.repository';
import { EntityModel } from '../../entities/EntityModel';

@Module({
  controllers: [TagsController],
  providers: [
    TagsRepository,
    TagsService,
    UserRepository,
    CollectionRepository,
    ElasticsearchService,
    EntityRepository,
  ],
  imports: [
    TypeOrmModule.forFeature([User, Collection, Tags, EntityModel]),
    ConfigModule.forRoot({
      envFilePath: '.env',
    }),
    JwtRegisterModule,
    HttpModule,
  ],
  exports: [TagsService],
})
export class TagsModule {}
