import { Controller, Post } from '@nestjs/common';
import { ElasticsearchService } from '../services/ElasticSearch.service';
import { Client } from '@elastic/elasticsearch';
import { InjectRepository } from '@nestjs/typeorm';
import { Collection } from '../entities/Collection';
import { Repository } from 'typeorm';
import { CollectionItem } from '../entities/CollectionItem';
import { Tags } from '../entities/Tags';
import { EntityModel } from '../entities/EntityModel';
import { User } from '../entities/User';

@Controller('recreate-elastic-search-index')
export class ReindexController {
  constructor(
    private readonly elasticsearchService: ElasticsearchService,
    @InjectRepository(Collection)
    private collectionRepository: Repository<Collection>,
    @InjectRepository(CollectionItem)
    private collectionItemRepository: Repository<CollectionItem>,
    @InjectRepository(Tags)
    private tagsRepository: Repository<Tags>,
    @InjectRepository(EntityModel)
    private entityRepository: Repository<EntityModel>,
    @InjectRepository(User)
    private userRepository: Repository<User>,
  ) {}

  @Post()
  async reindex() {
    const client = new Client({
      node: `http://${process.env.ELASTIC_URL}:9200`,
    });
    const prepareBulkData = (
      index: string,
      items: any[],
      mapper: (item: any) => any,
    ) => {
      return items.flatMap((item) => [
        { index: { _index: index, _id: item.id } },
        mapper(item),
      ]);
    };
    try {
      const users = await this.userRepository.find();
      const collections = await this.collectionRepository.find();
      const collectionItems = await this.collectionItemRepository.find();
      const tags = await this.tagsRepository.find();
      const entities = await this.entityRepository.find();

      const bulkActions = [
        ...prepareBulkData('user', users, (user) => ({
          login: user.login,
          email: user.email,
          user_name: user.user_name,
        })),
        ...prepareBulkData('collection', collections, (collection) => ({
          id: collection.id,
          name: collection.name,
        })),
        ...prepareBulkData(
          'collection_item',
          collectionItems,
          (collectionItem) => ({
            id: collectionItem.id,
            name: collectionItem.name,
          }),
        ),
        ...prepareBulkData('tags', tags, (tag) => ({
          id: tag.id,
          name: tag.name,
        })),
        ...prepareBulkData('entity_model', entities, (entity) => ({
          id: entity.id,
          name: entity.name,
        })),
      ];

      await client.bulk({
        body: bulkActions,
        refresh: true,
      });

      return { success: true, message: 'Reindexing completed' };
    } catch (error) {
      console.error('Reindexing failed:', error);
      throw new Error('Reindexing failed');
    }
    // return await this.elasticsearchService.reindexAllData();
  }
}
