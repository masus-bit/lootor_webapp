import { ApiProperty } from '@nestjs/swagger';
import {
  ReturnedCollectionDto,
  ReturnUser,
} from '../collections/CreateCollectionDto';
import { ReturnedCollectionItemDto } from '../collectionItem/CollectionItemCreateDto';
import { Event } from '../../entities/Event';

export class EventsDto {
  @ApiProperty()
  readonly id: string;
  @ApiProperty()
  readonly action: string;
  @ApiProperty()
  readonly eventTarget: string;
  @ApiProperty()
  readonly date: Date;
  @ApiProperty()
  readonly targetName: string;
  @ApiProperty({ type: ReturnUser })
  readonly user: string;
  @ApiProperty({ type: ReturnUser })
  readonly targetUser?: string;
  @ApiProperty({ type: ReturnedCollectionDto })
  readonly targetCollection?: string;
  @ApiProperty({ type: ReturnedCollectionItemDto })
  readonly targetCollectionItem?: string;

  constructor(event: Readonly<Event>) {
    this.id = event.id;
    this.date = event.date;
    this.action = event.action;
    this.eventTarget = event.event_target;
    this.targetName = event.target_name;
    // @ts-ignore
    this.user = new ReturnUser(event.user);
    // @ts-ignore
    this?.targetUser = event?.target_user
      ? // @ts-ignore
        new ReturnUser(event?.target_user)
      : null;
    // @ts-ignore
    this?.targetCollection = event?.target_collection
      ? // @ts-ignore
        new ReturnedCollectionDto(event?.target_collection)
      : null;
    // @ts-ignore
    this?.targetCollectionItem = event?.target_collection_item
      ? new ReturnedCollectionItemDto(
          // @ts-ignore
          event?.target_collection_item,
        )
      : null;
  }
}

export class GetEventsDto {
  @ApiProperty({ type: EventsDto, isArray: true })
  readonly data: ReadonlyArray<EventsDto>;

  constructor(data: ReadonlyArray<EventsDto>) {
    this.data = data;
  }
}
