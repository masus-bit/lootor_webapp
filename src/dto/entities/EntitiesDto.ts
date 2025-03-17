import { Exclude, Expose } from 'class-transformer';
import { ApiProperty } from '@nestjs/swagger';

@Exclude()
export class CreateEntitiesDto {
  @Expose()
  @ApiProperty()
  readonly names: string[];
}
export class CreateEntityDto {
  @Expose()
  @ApiProperty()
  readonly name: string;
}
