import { Exclude, Expose } from 'class-transformer';
import { ApiProperty } from '@nestjs/swagger';

@Exclude()
export class CreateEventDto {
  @Expose()
  @ApiProperty()
  readonly user: string;

  @Expose()
  @ApiProperty()
  readonly date: number;

  @Expose({ name: 'targetUser' })
  @ApiProperty({ name: 'targetUser' })
  readonly target_user: string;

  @Expose({ name: 'targetCollection' })
  @ApiProperty({ name: 'targetCollection' })
  readonly target_collection: string;

  @Expose({ name: 'targetCollectionItem' })
  @ApiProperty({ name: 'targetCollectionItem' })
  readonly target_collection_item: string;

  @Expose()
  @ApiProperty()
  readonly action: string;

  @Expose({ name: 'targetName' })
  @ApiProperty({ name: 'targetName' })
  readonly target_name: string;

  @Expose({ name: 'eventTarget' })
  @ApiProperty({ name: 'eventTarget' })
  readonly event_target: string;
}
