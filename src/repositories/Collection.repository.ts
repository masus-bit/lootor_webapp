import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Collection } from '../entities/Collection';
import { Repository } from 'typeorm';
import { DeepPartial } from 'typeorm/common/DeepPartial';
import { ElasticsearchService } from '../services/ElasticSearch.service';

@Injectable()
export class CollectionRepository {
  constructor(
    @InjectRepository(Collection)
    private collectionRepository: Repository<Collection>,
    private readonly elasticsearchService: ElasticsearchService,
  ) {}

  async find(): Promise<Collection[]> | never {
    return await this.collectionRepository.find();
  }

  async save(data: DeepPartial<Collection>): Promise<Collection> {
    const collection = await this.collectionRepository.save(data);
    await this.elasticsearchService.createIndexIfNotExists('collection');
    await this.elasticsearchService.upsertDocument(
      'collection',
      collection.id.toString(),
      {
        name: collection.name,
        id: collection.id,
      },
    );
    return collection;
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
        .leftJoinAndSelect('collection.tags', 'tags')
        .leftJoinAndSelect('collection.collectionItems', 'collection_item')
        // .leftJoinAndMapMany(
        //   'collection.collectionItems',
        //   CollectionItem,
        //   'collection_item',
        //   'collection_item.collection = collection.id AND collection_item.deleted = false',
        // )
        .getOne();
    } catch (e) {
      console.log(e.message);
    }
  }

  async getByShareString(shareString: string): Promise<Collection | never> {
    try {
      return await this.collectionRepository
        .createQueryBuilder('collection')
        .where('collection.share_string = :shareString', { shareString })
        .andWhere('collection.deleted = :deleted', { deleted: false })
        .leftJoinAndSelect('collection.user', 'user')
        .leftJoinAndSelect('collection.tags', 'tags')
        .leftJoinAndSelect('collection.collectionItems', 'collection_item')
        // .leftJoinAndMapMany(
        //   'collection.collectionItems',
        //   CollectionItem,
        //   'collection_item',
        //   'collection_item.collection = collection.id AND collection_item.deleted = false',
        // )
        .getOne();
    } catch (e) {
      console.log(e.message);
    }
  }

  async getByIdWithoutCollections(id: string): Promise<Collection | never> {
    return await this.collectionRepository
      .createQueryBuilder('collection')
      .leftJoinAndSelect('collection.user', 'user')
      .leftJoinAndSelect('collection.tags', 'tags')
      .where('collection.id = :id', { id })
      .andWhere('collection.deleted = :deleted', { deleted: false })
      .getOne();
  }

  async getOneByTransliteration(
    login: string,
    transliteration: string,
  ): Promise<Collection | null> {
    try {
      return await this.collectionRepository
        .createQueryBuilder('collection')
        .leftJoinAndSelect('collection.user', 'user')
        .leftJoinAndSelect('collection.tags', 'tags')
        .leftJoinAndSelect('collection.collectionItems', 'collection_item')
        .leftJoinAndSelect('collection_item.platform', 'platforms')
        .leftJoinAndSelect('collection_item.entities', 'entity_model')
        .where('LOWER(user.login) = LOWER(:login)', { login })
        .andWhere('collection.transliteration = :translit', {
          translit: transliteration,
        })
        .andWhere('collection.deleted = :deleted', { deleted: false })
        .getOne();
    } catch (err) {
      console.error(err);
      return null;
    }
  }

  async getByUserId(id: string): Promise<Collection[]> {
    return await this.collectionRepository
      .createQueryBuilder('collection')
      .leftJoinAndSelect('collection.user', 'user')
      .leftJoinAndSelect('collection.tags', 'tags')
      .leftJoinAndSelect('collection.collectionItems', 'collection_item')
      .leftJoinAndSelect('collection_item.platform', 'platforms')
      .leftJoinAndSelect('collection_item.entities', 'entity_model')
      .where('LOWER(user.login) = LOWER(:id)', { id })
      .andWhere('collection.deleted = :deleted', { deleted: false })
      .orderBy('collection.created', 'DESC')
      .getMany();
  }

  async getByUserIdWithoutItems(id: string): Promise<Collection[]> {
    return await this.collectionRepository
      .createQueryBuilder('collection')
      .leftJoinAndSelect('collection.user', 'user')
      .leftJoinAndSelect('collection.tags', 'tags')
      .where('LOWER(user.login) = LOWER(:id)', { id })
      .andWhere('collection.deleted = :deleted', { deleted: false })
      .orderBy('collection.created', 'DESC')
      .getMany();
  }

  async getByUserIdWithoutPrivates(id: string): Promise<Collection[]> {
    return await this.collectionRepository
      .createQueryBuilder('collection')
      .innerJoinAndSelect('collection.user', 'user')
      .leftJoinAndSelect('collection.tags', 'tags')
      .where('LOWER(user.login) = LOWER(:id)', { id })
      .andWhere('collection.deleted = :deleted', { deleted: false })
      .andWhere('collection.is_private = :isPrivate', { isPrivate: false })
      .orderBy('collection.created', 'DESC')
      .getMany();
  }

  async getByIdWithoutUser(id: string): Promise<Collection | never> {
    return await this.collectionRepository
      .createQueryBuilder('collection')
      .where('collection.id = :id', { id })
      .andWhere('collection.deleted = :deleted', { deleted: false })
      .leftJoinAndSelect('collection.tags', 'tags')
      .getOne();
  }

  async getByTag(tag: string): Promise<Collection[] | never> {
    return await this.collectionRepository
      .createQueryBuilder('collection')
      .leftJoinAndSelect('collection.user', 'user')
      .leftJoinAndSelect('collection.tags', 'tags')
      .leftJoinAndSelect('collection.collectionItems', 'collection_item')
      .leftJoinAndSelect('collection_item.platform', 'platforms')
      .leftJoinAndSelect('collection_item.entities', 'entity_model')
      .leftJoin('collection.tags', 'tagsForFilter')
      .where('collection.is_private = :isPrivate', { isPrivate: false })
      .andWhere('tagsForFilter.name = :tag', { tag })
      .orderBy('collection.created', 'ASC')
      .getMany();
  }

  async updateById(
    collection: DeepPartial<Collection>,
    id: string,
  ): Promise<Collection> {
    const updatedCollection = this.collectionRepository.create({
      id,
      ...collection,
    });
    await this.elasticsearchService.createIndexIfNotExists('collection');

    await this.elasticsearchService.indexDocument('collection', {
      name: updatedCollection.name,
    });
    console.log('INDEXED');
    return await this.collectionRepository.save(updatedCollection);
  }

  async getCountByUserLogin(login: string): Promise<number | never> {
    return await this.collectionRepository
      .createQueryBuilder('collection')
      .where('collection.user = :login', { login })
      .getCount();
  }

  async delete(id: string) {
    const del = await this.collectionRepository.delete({ id });
    await this.elasticsearchService.deleteDocument('collection', id);
    return del;
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
