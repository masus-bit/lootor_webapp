import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { ReindexController } from '../../controllers/Reindex.controller';
import { User } from '../../entities/User';
import { Collection } from '../../entities/Collection';
import { Tags } from '../../entities/Tags';
import { EntityModel } from '../../entities/EntityModel';
import { CollectionItem } from '../../entities/CollectionItem';
import { ElasticsearchService } from '../../services/ElasticSearch.service';

@Module({
  imports: [
    TypeOrmModule.forFeature([
      User,
      Collection,
      Tags,
      EntityModel,
      CollectionItem,
    ]),
  ],
  controllers: [ReindexController],
  providers: [ElasticsearchService],
})
export class ReindexModule {}
