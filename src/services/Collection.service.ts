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

      const result = await this.collectionRepository.getByIdWithoutCollections(
        collection.id,
      );
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
      const result = await this.collectionRepository.getByIdWithoutCollections(
        collection.id,
      );
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

  async getOne(id?: string, transliteration?: string, login?: string) {
    let collection;
    if (id) {
      collection = await this.collectionRepository.getById(id);
    } else {
      collection = await this.collectionRepository.getOneByTransliteration(
        login,
        transliteration,
      );
    }
    return new GetOneCollectionDto(new CollectionDto(collection));
  }

  async like(id: string, userId: string) {
    const exists = await this.collectionRepository.getByIdWithoutCollections(
      id,
    );

    const model = await this.collectionRepository.createModel(exists);
    if (exists?.likes?.length) {
      let likes = [];
      !model.likes.includes(userId)
        ? // @ts-ignore
          (model.likes = `{${exists.likes}, ${userId}}`)
        : (likes = model.likes.filter((l) => l !== userId));
      // @ts-ignore
      model.likes = `{${likes}}`;
      await this.collectionRepository.save(model);
      return 'Liked';
    }

    // @ts-ignore
    model.likes = `{${userId}}`;
    await this.collectionRepository.save(model);
    return 'Liked';
  }
}
