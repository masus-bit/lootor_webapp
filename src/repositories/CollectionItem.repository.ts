import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { DeepPartial } from 'typeorm/common/DeepPartial';
import { CollectionItem } from '../entities/CollectionItem';

@Injectable()
export class CollectionItemRepository {
  constructor(
    @InjectRepository(CollectionItem)
    private collectionItemRepository: Repository<CollectionItem>,
  ) {}

  async save(data: DeepPartial<CollectionItem>): Promise<CollectionItem> {
    return await this.collectionItemRepository.save(data);
  }

  createModel(data: DeepPartial<CollectionItem>): CollectionItem {
    return this.collectionItemRepository.create(data);
  }

  async getById(id: string): Promise<CollectionItem | never> {
    return await this.collectionItemRepository
      .createQueryBuilder('collection_item')
      .where('collection_item.id = :id', { id })
      .andWhere('collection_item.deleted = :deleted', { deleted: false })
      .innerJoinAndSelect('collection_item.collection', 'collection')
      .innerJoinAndSelect('collection.user', 'user')
      .getOne();
  }

  async updateById(
    collectionItem: DeepPartial<CollectionItem>,
    id: string,
  ): Promise<CollectionItem> {
    // @ts-ignore
    const updatedCollectionItem = this.collectionItemRepository.create({
      id,
      item_photos: `{${collectionItem.item_photos}}`,
      ...collectionItem,
    });
    // @ts-ignore
    return await this.collectionItemRepository.save(updatedCollectionItem);
  }

  async delete(id: string) {
    const collectionItem = await this.collectionItemRepository
      .createQueryBuilder('collection_item')
      .where('collection_item.id = :id', { id })
      .getOne();
    const updatedCollectionItem = this.collectionItemRepository.create({
      id,
      deleted: true,
      ...collectionItem,
    });
    await this.collectionItemRepository.save(updatedCollectionItem);
    return 'Deleted successfully';
  }
}
