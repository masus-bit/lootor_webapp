import { Inject, Injectable } from '@nestjs/common';
import { CollectionItemRepository } from '../repositories/CollectionItem.repository';
import {
  CollectionItemCreateDto,
  ReturnCreateCollectionItem,
  ReturnedCollectionItemDto,
} from '../dto/collectionItem/CollectionItemCreateDto';
import { HttpInternalServerError } from '../errors/HttpInternalServerError';

@Injectable()
export class CollectionItemService {
  constructor(
    @Inject(CollectionItemRepository)
    private collectionItemRepository: CollectionItemRepository,
  ) {}

  /**
   *
   * Создает запись collection item
   */
  async create(
    dto: CollectionItemCreateDto,
  ): Promise<ReturnCreateCollectionItem> {
    try {
      console.log(dto);
      const model = this.collectionItemRepository.createModel(dto);
      // @ts-ignore
      model.item_photos = `{${dto.item_photos}}`;
      const instance = await this.collectionItemRepository.save(model);
      const result = await this.collectionItemRepository.getById(instance.id);
      return new ReturnCreateCollectionItem(
        new ReturnedCollectionItemDto(result),
      );
    } catch (err) {
      console.error(err.message);
      throw err;
    }
  }

  /**
   * Удаляет запись коллекции
   */
  async delete(id: string) {
    const deleted = await this.collectionItemRepository.delete(id);

    if (deleted.affected > 0) {
      return 'Collection item has deleted successfully!';
    } else {
      throw new Error('Nothing to delete');
    }
  }

  /**
   * изменяет запись коллекции
   */
  async update(dto: CollectionItemCreateDto, id: string) {
    const exists = await this.collectionItemRepository.getById(id);
    if (exists) {
      const model = await this.collectionItemRepository.createModel(dto);
      // @ts-ignore
      model.item_photos = `{${dto.item_photos}}`;
      const instance = await this.collectionItemRepository.save({
        id,
        ...model,
      });
      const result = await this.collectionItemRepository.getById(instance.id);
      return new ReturnCreateCollectionItem(
        new ReturnedCollectionItemDto(result),
      );
    }
  }

  async getById(id: string) {
    try {
      const exists = await this.collectionItemRepository.getById(id);
      return new ReturnCreateCollectionItem(
        new ReturnedCollectionItemDto(exists),
      );
    } catch (err) {
      throw new HttpInternalServerError(
        'Что-то пошло не так, обратитесь к кому-либо',
      );
    }
  }
}
