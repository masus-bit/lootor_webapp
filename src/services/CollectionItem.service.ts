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
      model.images = `{${dto.images}}`;
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
          result.name,
          null,
          null,
          instance.id,
        );
      }
      return new ReturnCreateCollectionItem(
        new ReturnedCollectionItemDto(result),
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
    const deleted = await this.collectionItemRepository.delete(id);

    if (deleted.affected > 0) {
      // @ts-ignore
      if (!exists.collection.is_private) {
        await this.eventRepository.addEvent(
          // @ts-ignore
          exists.collection.user.login,
          EventActions.delete,
          EventTargets.collectionItem,
          exists.name,
        );
      }
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
      const data = { ...dto };
      delete data.collection;
      // const model = this.collectionItemRepository.createModel(data);
      // @ts-ignore
      // model.images = `{${dto.images}}`;
      // console.log({
      //   id,
      //   ...exists,
      //   ...data,
      //   images: `{${dto.images}}`,
      //   entities: [...exists.entities, ...dto.entities],
      // });
      try {
        let result = [];

        if (data?.entities) {
          result = await this.getEntities(data?.entities);
        }
        const instance = await this.collectionItemRepository.save({
          id,
          ...exists,
          ...data,
          // @ts-ignore
          images: `{${data.images}}`,
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
    console.log(id);
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
}
