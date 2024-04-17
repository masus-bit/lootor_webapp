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
import { getUnixDate } from '../utils/date';

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
      collectionModel.created = getUnixDate();
      collectionModel.share_string = randomstring.generate(8);
      const collection = await this.collectionRepository.save(collectionModel);
      collection.totalPrice = 0;
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
        return 'Collection has deleted successfully!';
      } else {
        throw new Error('Nothing to delete');
      }
    } catch (e) {
      console.log(e);
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
        { ...dto, tags: [...existsCollection.tags, ...dto.tags] },
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
  async getByUserId(
    id: string,
    authorizedUser?: string,
  ): Promise<ReturnCollectionsDto> {
    let collections;
    if (authorizedUser === id) {
      collections = await this.collectionRepository.getByUserId(id);
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
        result.push(
          new ReturnedCollectionDto(
            item,
            authorizedUser ? !subArray.includes(item.id) : false,
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
        authorizedUser ? !subArray.includes(collection.id) : false,
      ),
    );
  }

  async like(id: string, userId: string) {
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

      return 'Liked';
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
    return 'Liked';
  }

  async subscribe(
    subscriptionTargetId: string,
    subscriber: User,
    isSubscribe: string,
  ): Promise<string> {
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
        return 'Подписка на коллекцию оформлена';
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
        return 'Отписка оформлена :D';
      }
    } catch (err) {
      console.log(err);
      throw new HttpBadRequestError('Что-то пошло не так');
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
