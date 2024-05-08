import { Inject, Injectable } from '@nestjs/common';
import { User } from '../entities/User';
import { EventRepository } from '../repositories/Event.repository';
import { EventsDto, GetEventsDto } from '../dto/events/GetEventsDto';

@Injectable()
export class EventService {
  constructor(
    @Inject(EventRepository)
    private eventRepository: EventRepository,
  ) {}

  async getEvents(authUser: User): Promise<any> {
    try {
      const subscriptions = authUser.subscriptions;
      if (authUser.subscriptions?.length) {
        const events = await this.eventRepository.getEvents(subscriptions);
        const result = [];
        events.map((e) => result.push(new EventsDto(e)));
        return new GetEventsDto(result);
      }
      return { data: [] };
    } catch (e) {
      console.log(e);
    }
  }

  async getFilteredEvents(
    userLogin: string | null,
    collectionId: string | null,
    collectionItem: string | null,
  ): Promise<any> {
    try {
      return await this.eventRepository.getFilteredEvents(
        userLogin,
        collectionId,
        collectionItem,
      );
      // const result = [];
      // events.map((e) => result.push(new EventsDto(e)));
      // return new GetEventsDto(result);
      //
      // return { data: [] };
    } catch (e) {
      console.log(e);
    }
  }
}
