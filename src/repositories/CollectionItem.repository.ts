import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { DeepPartial } from 'typeorm/common/DeepPartial';
import { CollectionItem } from '../entities/CollectionItem';
import { ElasticsearchService } from '../services/ElasticSearch.service';

@Injectable()
export class CollectionItemRepository {
  constructor(
    @InjectRepository(CollectionItem)
    private collectionItemRepository: Repository<CollectionItem>,
    private readonly elasticsearchService: ElasticsearchService,
  ) {}

  async save(data: DeepPartial<CollectionItem>): Promise<CollectionItem> {
    const collectionItem = await this.collectionItemRepository.save(data);

    await this.elasticsearchService.createIndexIfNotExists('collection_item', {
      properties: {
        name: {
          type: 'text',
          analyzer: 'common_analyzer',
        },
      },
    });
    await this.elasticsearchService.indexDocument('collection_item', {
      id: collectionItem.id,
      name: collectionItem.name,
    });
    return collectionItem;
  }

  createModel(data: DeepPartial<CollectionItem>): CollectionItem {
    return this.collectionItemRepository.create(data);
  }

  async getById(id: string): Promise<CollectionItem | never> {
    return await this.collectionItemRepository
      .createQueryBuilder('collection_item')
      .where('collection_item.id = :id', { id })
      .andWhere('collection_item.deleted = :deleted', { deleted: false })
      .leftJoinAndSelect('collection_item.collections', 'collection')
      .leftJoinAndSelect('collection.user', 'user')
      .leftJoinAndSelect('collection_item.platform', 'platforms')
      .leftJoinAndSelect('collection_item.entities', 'entity_model')
      .getOne();
  }

  async getCountByUserLogin(login: string): Promise<number | never> {
    return await this.collectionItemRepository
      .createQueryBuilder('collection_item')
      .where('collection_item.owner = :login', { login })
      .andWhere('collection_item.deleted = :deleted', { deleted: false })
      .getCount();
  }

  async updateById(
    collectionItem: DeepPartial<CollectionItem>,
    id: string,
  ): Promise<CollectionItem> {
    // @ts-ignore
    const updatedCollectionItem = this.collectionItemRepository.create({
      id,
      images: `{${collectionItem.images}}`,
      ...collectionItem,
    });
    // @ts-ignore
    return await this.collectionItemRepository.save(updatedCollectionItem);
  }

  // sum(collectionId: string) {
  //   return this.collectionItemRepository.sum('purchase_price', {
  //     //@ts-ignore
  //     collections: collectionId,
  //   });
  // }

  async sum(collectionId: string): Promise<number> {
    const result = await this.collectionItemRepository
      .createQueryBuilder('item')
      .select('SUM(item.purchase_price)', 'sum')
      .innerJoin(
        'item.collections',
        'collection',
        'collection.id = :collectionId',
        { collectionId },
      )
      .where('item.deleted = :deleted', { deleted: false })
      .getRawOne();

    return result?.sum || 0;
  }

  async sumShippingCost(collectionId: string): Promise<number> {
    const result = await this.collectionItemRepository
      .createQueryBuilder('item')
      .select('SUM(item.shipping_cost)', 'sum')
      .innerJoin(
        'item.collections',
        'collection',
        'collection.id = :collectionId',
        { collectionId },
      )
      .where('item.deleted = :deleted', { deleted: false })
      .getRawOne();

    return result?.sum || 0;
  }

  sumByUserLogin(userLogin: string) {
    return this.collectionItemRepository.sum('purchase_price', {
      owner: userLogin,
    });
  }

  async delete(id: string) {
    const deleted = await this.collectionItemRepository.delete({ id });
    await this.elasticsearchService.deleteDocument(
      'collection_item',
      id.toString(),
    );
    return deleted;

    // const collectionItem = await this.collectionItemRepository
    //   .createQueryBuilder('collection_item')
    //   .where('collection_item.id = :id', { id })
    //   .getOne();
    // const updatedCollectionItem = this.collectionItemRepository.create({
    //   id,
    //   deleted: true,
    //   ...collectionItem,
    // });
    // await this.collectionItemRepository.save(updatedCollectionItem);
    // return 'Deleted successfully';
  }
}
