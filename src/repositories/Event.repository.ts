import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { DeepPartial } from 'typeorm/common/DeepPartial';
import { plainToClass } from 'class-transformer';
import { CreateEventDto } from '../dto/events/CreateEventDto';
import { Event } from '../entities/Event';
import { CollectionItem } from '../entities/CollectionItem';
import { User } from '../entities/User';
import { Collection } from '../entities/Collection';
import { getUnixDate } from '../utils/date';

@Injectable()
export class EventRepository {
  constructor(
    @InjectRepository(Event)
    private eventRepository: Repository<Event>,
  ) {}

  async save(data: DeepPartial<Event>): Promise<Event> {
    return await this.eventRepository.save(data);
  }

  createModel(data: CreateEventDto): Event {
    return this.eventRepository.create(data);
  }

  async addEvent(
    user: string,
    action: string,
    eventTarget: string,
    targetName: string,
    targetUser?: string,
    targetCollection?: string,
    targetCollectionItem?: string,
  ): Promise<Event> {
    const model = this.createModel(
      plainToClass(CreateEventDto, {
        user,
        action,
        eventTarget,
        targetName,
        targetUser,
        targetCollection,
        targetCollectionItem,
        date: getUnixDate(),
      }),
    );
    return await this.save(model);
  }

  async getEvents(subscriptions: string[]): Promise<Event[] | never> {
    try {
      return await this.eventRepository
        .createQueryBuilder('event')
        .where('event.user IN (:...subscriptions)', { subscriptions })
        .leftJoinAndSelect('event.user', 'user as other')
        .leftJoinAndMapOne(
          'event.target_user',
          User,
          'user',
          'user.login = event.target_user',
        )
        .leftJoinAndMapOne(
          'event.target_collection',
          Collection,
          'collection',
          'collection.id = event.target_collection',
        )
        .leftJoinAndMapOne(
          'event.target_collection_item',
          CollectionItem,
          'collection_item',
          'collection_item.id = event.target_collection_item',
        )
        .orderBy('date', 'DESC')
        .getMany();
    } catch (e) {
      console.log(e);
    }
  }

  async getFilteredEvents(
    userLogin: string | null,
    collectionId: string | null,
    collectionItemId: string | null,
  ): Promise<Event[] | never> {
    try {
      return await this.eventRepository
        .createQueryBuilder('event')
        .where('event.user = :userLogin', { userLogin })
        .orWhere('event.target_collection = :collectionId', { collectionId })
        .orWhere('event.target_collection_item = :collectionItemId', {
          collectionItemId,
        })
        .leftJoinAndMapOne(
          'event.target_user',
          User,
          'user',
          'user.login = event.target_user',
        )
        .leftJoinAndMapOne(
          'event.target_collection',
          Collection,
          'collection',
          'collection.id = event.target_collection',
        )
        .leftJoinAndMapOne(
          'event.target_collection_item',
          CollectionItem,
          'collection_item',
          'collection_item.id = event.target_collection_item',
        )
        .orderBy('date', 'DESC')
        .getMany();
    } catch (e) {
      console.log(e);
    }
  }

  async deleteEvent(
    user: string,
    targetName: string,
    eventTarget,
    target: string,
  ): Promise<any> {
    const targetField = {
      targetUser: 'target_user',
      targetCollection: 'target_collection',
      targetCollectionItem: 'target_collection_item',
    };
    const exists = await this.eventRepository
      .createQueryBuilder('event')
      .where('event.user =:user', { user })
      .andWhere('event.target_name =:targetName', { targetName })
      .andWhere('event.event_target =:eventTarget', { eventTarget })
      .andWhere(`event.${targetField[eventTarget]} =:target`, { target })
      .getOne();
    if (exists) {
      await this.eventRepository.delete({ id: exists.id });
    }
  }
}
