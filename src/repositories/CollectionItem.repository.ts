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
      .innerJoinAndSelect('collection_item.platform', 'platforms')
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

  sum(collectionId: string) {
    return this.collectionItemRepository.sum('purchase_price', {
      collection: collectionId,
    });
  }

  sumShippingCost(collectionId: string) {
    return this.collectionItemRepository.sum('shipping_cost', {
      collection: collectionId,
    });
  }

  sumByUserLogin(userLogin: string) {
    return this.collectionItemRepository.sum('purchase_price', {
      owner: userLogin,
    });
  }

  async delete(id: string) {
    return await this.collectionItemRepository.delete({ id });
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
