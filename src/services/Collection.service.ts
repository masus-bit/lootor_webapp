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
import { EventRepository } from '../repositories/Event.repository';
import { EventActions, EventTargets } from '../types/base';
import { User } from '../entities/User';
import { HttpBadRequestError } from '../errors/HttpBadRequestError';
import { UserRepository } from '../repositories/User.repository';
import { CollectionItemRepository } from '../repositories/CollectionItem.repository';
import { defineShareString } from '../utils/defineShareString';
import { getIsoDate } from '../utils/date';
import { TagsRepository } from '../repositories/Tags.repository';
import { ReturnDataDto } from '../dto/ReturnDataDto';

var randomstring = require('randomstring');

@Injectable()
export class CollectionService {
  constructor(
    @Inject(CollectionRepository)
    private collectionRepository: CollectionRepository,
    @Inject(EventRepository)
    private eventRepository: EventRepository,
    @Inject(UserRepository)
    private userRepository: UserRepository,
    @Inject(CollectionItemRepository)
    private collectionItemRepository: CollectionItemRepository,
    @Inject(TagsRepository)
    private tagsRepository: TagsRepository,
  ) {}

  /**
   *
   * Создает запись коллекции
   */
  async createCollection(
    dto: CreateCollectionDto,
  ): Promise<ReturnCreateCollection> {
    try {
      const resultTags = [];
      for (const tag of dto.tags) {
        const exist = await this.tagsRepository.getTagByName(tag);
        if (!exist) {
          const saved = await this.tagsRepository.addTag(tag);
          resultTags.push(saved);
        } else {
          resultTags.push(exist);
        }
      }
      const collectionModel = this.collectionRepository.createModel({
        ...dto,
        tags: resultTags,
      });
      collectionModel.created = getIsoDate();
      collectionModel.share_string = randomstring.generate(8);
      const collection = await this.collectionRepository.save(collectionModel);
      collection.totalPrice = 0;
      collection.shippingTotal = 0;
      const result = await this.collectionRepository.getByIdWithoutCollections(
        collection.id,
      );
      if (!result.is_private) {
        await this.eventRepository.addEvent(
          // @ts-ignore
          result.user.login,
          EventActions.create,
          EventTargets.collection,
          result.name,
          null,
          collection.id,
        );
      }
      return new ReturnCreateCollection(
        new ReturnedCollectionDto(result, true),
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
    try {
      const exists = await this.collectionRepository.getByIdWithoutCollections(
        id,
      );
      const deleted = await this.collectionRepository.delete(id);

      if (deleted.affected > 0) {
        if (!exists.is_private) {
          await this.eventRepository.addEvent(
            // @ts-ignore
            exists.user.login,
            EventActions.delete,
            EventTargets.collection,
            exists.name,
          );
        }
        return { data: { success: true } };
      } else {
        throw new Error('Nothing to delete');
      }
    } catch (e) {
      throw new HttpBadRequestError('Nothing to delete');
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
      const resultTags = [];
      for (const tag of dto.tags) {
        const exist = await this.tagsRepository.getTagByName(tag);
        if (!exist) {
          const saved = await this.tagsRepository.addTag(tag);
          resultTags.push(saved);
        } else {
          resultTags.push(exist);
        }
      }
      const collection = await this.collectionRepository.updateById(
        { ...dto, tags: [...existsCollection.tags, ...resultTags] },
        collectionId,
      );
      const result = await this.collectionRepository.getById(collection.id);
      return new ReturnCreateCollection(new ReturnedCollectionDto(result));
    }
  }

  /**
   * ПОлучение всех коллекциq по логину юзера
   */
  async getByUserId(
    id: string,
    authorizedUser?: string,
  ): Promise<ReturnCollectionsDto> {
    let collections;
    if (authorizedUser?.toLowerCase() === id.toLowerCase()) {
      collections = await this.collectionRepository.getByUserIdWithoutItems(id);
    } else {
      collections = await this.collectionRepository.getByUserIdWithoutPrivates(
        id,
      );
    }
    let authUser;
    let subArray;
    if (authorizedUser) {
      authUser = await this.userRepository.getById(authorizedUser);
      subArray = authUser.collection_subscriptions;
    }
    const result: ReturnedCollectionDto[] = [];

    const process = async (array) => {
      for (const item of array) {
        item.share_string = defineShareString(authorizedUser, id, item);
        item.totalPrice =
          (await this.collectionItemRepository.sum(item.id)) || 0;
        item.shippingTotal =
          (await this.collectionItemRepository.sumShippingCost(item.id)) || 0;
        const count = await this.collectionItemRepository.getCount(item.id);
        result.push(
          new ReturnedCollectionDto(
            item,
            authorizedUser ? !subArray.includes(item.id) : true,
            count,
          ),
        );
      }
    };
    await process(collections);
    // collections.forEach((collection) => {
    //   collection.totalPrice = this.collectionItemRepository.sum(collection.id);
    //   result.push(new ReturnedCollectionDto(collection));
    // });
    return new ReturnCollectionsDto(result);
  }

  async getOne(
    authorizedUser?: string,
    id?: string,
    transliteration?: string,
    login?: string,
    shareString?: string,
  ) {
    let collection;
    if (id) {
      collection = await this.collectionRepository.getById(id);
      collection.totalPrice = await this.collectionItemRepository.sum(
        collection.id,
      );
      collection.shippingTotal =
        await this.collectionItemRepository.sumShippingCost(collection.id);
      collection.share_string = defineShareString(
        authorizedUser,
        collection.user.login,
        collection,
      );
    } else if (login) {
      collection = await this.collectionRepository.getOneByTransliteration(
        login,
        transliteration,
      );
      collection.totalPrice = await this.collectionItemRepository.sum(
        collection.id,
      );
      collection.shippingTotal =
        await this.collectionItemRepository.sumShippingCost(collection.id);
      collection.share_string = defineShareString(
        authorizedUser,
        collection.user.login,
        collection,
      );
    } else if (shareString) {
      collection = await this.collectionRepository.getByShareString(
        shareString,
      );

      collection.totalPrice = await this.collectionItemRepository.sum(
        collection.id,
      );
      collection.shippingTotal =
        await this.collectionItemRepository.sumShippingCost(collection.id);
      collection.share_string = '';
    }

    let authUser;
    let subArray;
    if (authorizedUser) {
      authUser = await this.userRepository.getById(authorizedUser);
      subArray = authUser.collection_subscriptions;
    }
    return new GetOneCollectionDto(
      new CollectionDto(
        collection,
        authorizedUser ? !subArray.includes(collection?.id) : true,
      ),
    );
  }

  async like(id: string, userId: string): Promise<ReturnDataDto> {
    const exists = await this.collectionRepository.getByIdWithoutCollections(
      id,
    );

    const model = await this.collectionRepository.createModel(exists);
    if (exists?.likes?.length) {
      let likes = [];
      if (!model.likes.includes(userId)) {
        // @ts-ignore
        model.likes = `{${exists.likes}, ${userId}}`;
        await this.eventRepository.addEvent(
          // @ts-ignore
          userId,
          EventActions.like,
          EventTargets.collection,
          exists.name,
          null,
          id,
        );
      } else {
        likes = model.likes.filter((l) => l !== userId);
        await this.eventRepository.deleteEvent(
          userId,
          exists.name,
          EventTargets.collection,
          exists.id,
        );
      }
      // @ts-ignore
      model.likes = `{${likes}}`;
      await this.collectionRepository.save(model);

      return new ReturnDataDto({
        success: true,
      });
    }

    // @ts-ignore
    model.likes = `{${userId}}`;
    await this.collectionRepository.save(model);
    await this.eventRepository.addEvent(
      // @ts-ignore
      userId,
      EventActions.like,
      EventTargets.collection,
      exists.name,
      null,
      exists.id,
    );
    return new ReturnDataDto({
      success: true,
    });
  }

  async subscribe(
    subscriptionTargetId: string,
    subscriber: User,
    isSubscribe: string,
  ): Promise<ReturnDataDto> {
    try {
      const subscriberModel = await this.userRepository.createModel(subscriber);
      const targetCollection = await this.collectionRepository.getById(
        subscriptionTargetId,
      );

      if (JSON.parse(isSubscribe)) {
        subscriberModel.collection_subscriptions?.length
          ? // @ts-ignore
            (subscriberModel.collection_subscriptions = `{ ${subscriberModel.collection_subscriptions}, ${subscriptionTargetId} }`)
          : // @ts-ignore
            (subscriberModel.collection_subscriptions = `{ ${subscriptionTargetId} }`);
        await this.userRepository.save(subscriberModel);
        await this.collectionRepository.updateById(
          {
            ...targetCollection,
            subscribers_count: (targetCollection.subscribers_count += 1),
          },
          subscriptionTargetId,
        );
        // TODO реализовать собтие для ленты "подписка на коллекцию" и сделать вывод в ленте
        // await this.eventRepository.addEvent(
        //   // @ts-ignore
        //   subscriber.login,
        //   EventActions.subscribe,
        //   EventTargets.user,
        //   subscriptionTargetId,
        //   subscriptionTargetId,
        // );
        return new ReturnDataDto({
          success: true,
        });
      } else {
        const subscribers = subscriberModel.collection_subscriptions.filter(
          (s) => s !== subscriptionTargetId,
        );
        // @ts-ignore
        subscriberModel.collection_subscriptions = `{${subscribers}}`;
        await this.userRepository.save(subscriberModel);
        // await this.eventRepository.deleteEvent(
        //   subscriber.login,
        //   subscriptionTargetId,
        //   EventTargets.user,
        //   subscriptionTargetId,
        // );
        await this.collectionRepository.updateById(
          {
            ...targetCollection,
            subscribers_count: (targetCollection.subscribers_count -= 1),
          },
          subscriptionTargetId,
        );
        return new ReturnDataDto({
          success: true,
        });
      }
    } catch (err) {
      console.log(err);
      throw new HttpBadRequestError();
    }
  }

  async getByTag(
    tag: string,
    authorizedUser: string,
  ): Promise<ReturnCollectionsDto> {
    const final = [];
    const result = await this.collectionRepository.getByTag(tag);
    let authUser;
    let subArray;
    if (authorizedUser) {
      authUser = await this.userRepository.getById(authorizedUser);
      subArray = authUser.collection_subscriptions;
    }
    result.map((collection) => {
      collection.share_string = defineShareString(
        authorizedUser,
        // @ts-ignore
        collection.user.login,
        collection,
      );

      final.push(
        new ReturnedCollectionDto(
          collection,
          authorizedUser ? !subArray.includes(collection.id) : false,
        ),
      );
    });

    return new ReturnCollectionsDto(final);
  }
}
