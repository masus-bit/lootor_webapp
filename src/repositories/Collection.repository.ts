import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Collection } from '../entities/Collection';
import { Repository } from 'typeorm';
import { DeepPartial } from 'typeorm/common/DeepPartial';
import { CollectionItem } from '../entities/CollectionItem';

@Injectable()
export class CollectionRepository {
  constructor(
    @InjectRepository(Collection)
    private collectionRepository: Repository<Collection>,
  ) {}

  async save(data: DeepPartial<Collection>): Promise<Collection> {
    return await this.collectionRepository.save(data);
  }

  createModel(data: DeepPartial<Collection>): Collection {
    return this.collectionRepository.create(data);
  }

  async getById(id: string): Promise<Collection | never> {
    try {
      return await this.collectionRepository
        .createQueryBuilder('collection')
        .where('collection.id = :id', { id })
        .andWhere('collection.deleted = :deleted', { deleted: false })
        .leftJoinAndSelect('collection.user', 'user')
        .leftJoinAndMapMany(
          'collection.collectionItems',
          CollectionItem,
          'collection_item',
          'collection_item.collection = collection.id AND collection_item.deleted = false',
        )
        .getOne();
    } catch (e) {
      console.log(e.message);
    }
  }

  async getByIdWithoutCollections(id: string): Promise<Collection | never> {
    return await this.collectionRepository
      .createQueryBuilder('collection')
      .where('collection.id = :id', { id })
      .andWhere('collection.deleted = :deleted', { deleted: false })
      .leftJoinAndSelect('collection.user', 'user')
      .getOne();
  }

  async getOneByTransliteration(
    login: string,
    transliteration: string,
  ): Promise<any> {
    try {
      return await this.collectionRepository
        .createQueryBuilder('collection')
        .where('collection.user = :login', { login })
        .andWhere('collection.transliteration = :translit', {
          translit: transliteration,
        })
        .andWhere('collection.deleted = :deleted', { deleted: false })
        .leftJoinAndSelect('collection.user', 'user')
        .leftJoinAndMapMany(
          'collection.collectionItems',
          CollectionItem,
          'collection_item',
          'collection_item.collection = collection.id',
        )
        .where('collection_item.deleted = :deleted', { deleted: false })
        .getOne();
    } catch (err) {
      console.log(err.message);
    }
  }

  async getByUserId(id: string): Promise<Collection[]> {
    return await this.collectionRepository
      .createQueryBuilder('collection')
      .where('collection.user = :id', { id })
      .andWhere('collection.deleted = :deleted', { deleted: false })
      .innerJoinAndSelect('collection.user', 'user')
      .getMany();
  }

  async getByIdWithoutUser(id: string): Promise<Collection | never> {
    return await this.collectionRepository
      .createQueryBuilder('collection')
      .where('collection.id = :id', { id })
      .andWhere('collection.deleted = :deleted', { deleted: false })
      .getOne();
  }

  async updateById(
    collection: DeepPartial<Collection>,
    id: string,
  ): Promise<Collection> {
    const updatedCollection = this.collectionRepository.create({
      id,
      ...collection,
    });
    return await this.collectionRepository.save(updatedCollection);
  }

  async delete(id: string) {
    return await this.collectionRepository.delete({ id });
    // const collection = await this.collectionRepository
    //   .createQueryBuilder('collection')
    //   .where('collection.id = :id', { id })
    //   .getOne();
    // const updatedCollection = this.collectionRepository.create({
    //   id,
    //   deleted: true,
    //   ...collection,
    // });
    // await this.collectionRepository.save(updatedCollection);
    // return 'Deleted successfully';
  }
}
