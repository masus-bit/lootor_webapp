import { Inject, Injectable } from '@nestjs/common';
import {
  CreateCollectionDto,
  ReturnCollectionsDto,
  ReturnCreateCollection,
  ReturnedCollectionDto,
} from '../dto/collections/CreateCollectionDto';
import { CollectionRepository } from '../repositories/Collection.repository';
import {
  CollectionDto,
  GetOneCollectionDto,
} from '../dto/collections/CollectionDto';

@Injectable()
export class CollectionService {
  constructor(
    @Inject(CollectionRepository)
    private collectionRepository: CollectionRepository,
  ) {}

  /**
   *
   * Создает запись коллекции
   */
  async createCollection(
    dto: CreateCollectionDto,
  ): Promise<ReturnCreateCollection> {
    try {
      const collectionModel = this.collectionRepository.createModel(dto);
      const collection = await this.collectionRepository.save(collectionModel);
      const result = await this.collectionRepository.getById(collection.id);
      return new ReturnCreateCollection(new ReturnedCollectionDto(result));
    } catch (err) {
      console.error(err.message);
      throw err;
    }
  }

  /**
   * Удаляет запись коллекции
   */
  async delete(id: string) {
    const deletedCollections = await this.collectionRepository.delete(id);

    if (deletedCollections.affected > 0) {
      return 'Collection has deleted successfully!';
    } else {
      throw new Error('Nothing to delete');
    }
  }

  /**
   * изменяет запись коллекции
   */
  async update(dto: CreateCollectionDto, collectionId: string) {
    const existsCollection = await this.collectionRepository.getByIdWithoutUser(
      collectionId,
    );
    if (existsCollection) {
      const collection = await this.collectionRepository.updateById(
        dto,
        collectionId,
      );
      const result = await this.collectionRepository.getById(collection.id);
      return new ReturnCreateCollection(new ReturnedCollectionDto(result));
    }
  }

  /**
   * ПОлучение всех коллекциq по логину юзера
   */
  async getByUserId(id: string): Promise<ReturnCollectionsDto> {
    const collections = await this.collectionRepository.getByUserId(id);
    const result: ReturnedCollectionDto[] = [];
    collections.forEach((collection) =>
      result.push(new ReturnedCollectionDto(collection)),
    );
    return new ReturnCollectionsDto(result);
  }

  async getOne(id: string) {
    const collection = await this.collectionRepository.getById(id);
    return new GetOneCollectionDto(new CollectionDto(collection));
  }
}
