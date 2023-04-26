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
      const events = await this.eventRepository.getEvents(subscriptions);
      const result = [];
      events.map((e) => result.push(new EventsDto(e)));
      return new GetEventsDto(result);
    } catch (e) {
      console.log(e);
    }
  }
}
