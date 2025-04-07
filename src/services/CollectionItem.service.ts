import { Inject, Injectable } from '@nestjs/common';
import { CollectionItemRepository } from '../repositories/CollectionItem.repository';
import {
  CollectionItemCreateDto,
  ReturnCreateCollectionItem,
  ReturnedCollectionItemDto,
} from '../dto/collectionItem/CollectionItemCreateDto';
import { HttpInternalServerError } from '../errors/HttpInternalServerError';
import { EventActions, EventTargets } from '../types/base';
import { EventRepository } from '../repositories/Event.repository';
import { User } from '../entities/User';
import { CollectionRepository } from '../repositories/Collection.repository';
import { EntityRepository } from '../repositories/Entity.repository';
import {
  ReturnCreateCollection,
  ReturnedCollectionDto,
} from '../dto/collections/CreateCollectionDto';
import { HttpBadRequestError } from '../errors/HttpBadRequestError';
import { ReturnDataDto } from '../dto/ReturnDataDto';

@Injectable()
export class CollectionItemService {
  constructor(
    @Inject(CollectionItemRepository)
    private collectionItemRepository: CollectionItemRepository,
    @Inject(EventRepository)
    private eventRepository: EventRepository,
    @Inject(CollectionRepository)
    private collectionRepository: CollectionRepository,
    @Inject(EntityRepository)
    private entityRepository: EntityRepository,
  ) {}

  /**
   *
   * Создает запись collection item
   */
  async create(
    dto: CollectionItemCreateDto,
    user: User,
  ): Promise<ReturnCreateCollectionItem> {
    try {
      const data = { ...dto };
      delete data.collection;
      const model = this.collectionItemRepository.createModel(data);
      const collection = await this.collectionRepository.getById(
        //@ts-ignore
        dto.collection,
      );
      model.collections = [collection];
      model.owner = user.login;
      // @ts-ignore
      if (dto?.images) model.images = `{${dto?.images}}`;
      // @ts-ignore
      if (dto?.copy_number) model.copy_number = `{${dto?.copy_number}}`;

      let resultEntities = [];
      if (data?.entities) {
        resultEntities = await this.getEntities(data?.entities);
      }
      const instance = await this.collectionItemRepository.save({
        ...model,
        entities: !dto?.entities
          ? model.entities
          : [...model?.entities, ...resultEntities],
      });
      const result = await this.collectionItemRepository.getById(instance.id);
      // @ts-ignore
      if (!collection.is_private) {
        await this.eventRepository.addEvent(
          // @ts-ignore
          collection.user.login,
          EventActions.create,
          EventTargets.collectionItem,
          instance.name,
          null,
          null,
          instance.id,
        );
      }
      return new ReturnCreateCollectionItem(
        new ReturnedCollectionItemDto(instance),
      );
    } catch (err) {
      console.error(err);
      throw err;
    }
  }

  /**
   * Удаляет запись коллекции
   */
  async delete(id: string) {
    const exists = await this.collectionItemRepository.getById(id);
    // @ts-ignore
    const deleted = await this.collectionItemRepository.delete(id);

    if (deleted.affected > 0) {
      // @ts-ignore
      if (!exists.collections[0].is_private) {
        await this.eventRepository.addEvent(
          // @ts-ignore
          exists.collections[0].user.login,
          EventActions.delete,
          EventTargets.collectionItem,
          exists.name,
        );
      }
      return new ReturnDataDto({
        success: true,
      });
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
      const data = { ...dto };
      delete data.collection;
      try {
        let result = [];
        if (data?.entities) {
          result = await this.getEntities(data?.entities);
        }
        const getImages = () => {
          if (exists?.images) {
            return data?.images
              ? `{${exists?.images}, ${data?.images}}`
              : `{${exists?.images || null}}`;
          }
          return `{${data?.images || null}}`;
        };
        const instance = await this.collectionItemRepository.save({
          id,
          ...exists,
          // @ts-ignore
          images: getImages(),
          // @ts-ignore
          copy_number: `{${data?.copy_number || exists?.copy_number || null} }`,
          entities: [...exists?.entities, ...result],
        });
        const resultCollectionItem =
          await this.collectionItemRepository.getById(instance.id);
        return new ReturnCreateCollectionItem(
          new ReturnedCollectionItemDto(resultCollectionItem),
        );
      } catch (e) {
        console.log(e);
      }
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

  async getEntities(array) {
    const result = [];
    for (const item of array) {
      result.push(await this.entityRepository.getEntityByName(item));
    }
    return result;
  }

  async copyOrMove(
    id: string,
    targetCollectionIds: string[],
    sourceCollectionId: string,
  ) {
    try {
      if (!sourceCollectionId) {
        const collections = [];
        for (const col of targetCollectionIds) {
          const collection = await this.collectionRepository.getById(col);
          collections.push(collection);
        }
        // const collection = await this.collectionRepository.getById(
        //   targetCollectionId,
        // );
        const collectionItem = await this.collectionItemRepository.getById(id);

        for (const col of collections) {
          try {
            await this.collectionRepository.updateById(
              {
                ...col,
                collectionItems: [...col.collectionItems, collectionItem],
              },
              col.id,
            );
          } catch (e) {
            throw new Error(e);
          }
        }
        // const result = await this.collectionRepository.getById(collection.id);
        // return new ReturnCreateCollection(new ReturnedCollectionDto(result));
        return new ReturnDataDto({
          success: true,
        });
      }
      if (targetCollectionIds?.length > 1)
        throw new HttpBadRequestError(
          'Перемещение экземпляра в несколько коллекций не поддерживается',
        );
      const sourceCollection = await this.collectionRepository.getById(
        sourceCollectionId,
      );
      const targetCollection = await this.collectionRepository.getById(
        targetCollectionIds[0],
      );
      const collectionItem = await this.collectionItemRepository.getById(id);
      await this.collectionRepository.updateById(
        {
          ...sourceCollection,
          collectionItems: sourceCollection.collectionItems.filter(
            (item) => item.id !== id,
          ),
        },
        sourceCollectionId,
      );
      await this.collectionRepository.updateById(
        {
          ...targetCollection,
          collectionItems: [
            ...targetCollection.collectionItems,
            collectionItem,
          ],
        },
        targetCollectionIds[0],
      );
      const result = await this.collectionRepository.getById(
        targetCollection.id,
      );
      return new ReturnCreateCollection(new ReturnedCollectionDto(result));
    } catch (err) {
      throw new HttpInternalServerError(err);
    }
  }
}
