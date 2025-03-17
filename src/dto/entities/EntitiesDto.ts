import { Exclude, Expose } from 'class-transformer';
import { ApiProperty } from '@nestjs/swagger';

@Exclude()
export class CreateEntitiesDto {
  @Expose()
  @ApiProperty()
  readonly names: string[];

  @Expose({ name: 'collectionItem' })
  @ApiProperty({ name: 'collectionItem' })
  readonly collection_item: string;
}
export class CreateEntityDto {
  @Expose()
  @ApiProperty()
  readonly name: string;

  @Expose({ name: 'collectionItem' })
  @ApiProperty({ name: 'collectionItem' })
  readonly collection_item: string;
}
