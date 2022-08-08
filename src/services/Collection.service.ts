import { Inject, Injectable } from '@nestjs/common';
import {
  CreateCollectionDto,
  ReturnCreateCollection,
  ReturnedCollectionDto,
} from '../dto/collections/CreateCollectionDto';
import { CollectionRepository } from '../repositories/Collection.repository';

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
    console.log(existsCollection);
    if (existsCollection) {
      const collection = await this.collectionRepository.updateById(
        dto,
        collectionId,
      );
      const result = await this.collectionRepository.getById(collection.id);
      return new ReturnCreateCollection(new ReturnedCollectionDto(result));
    }
  }
}
