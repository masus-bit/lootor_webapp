import { Exclude, Expose } from 'class-transformer';
import { ApiProperty } from '@nestjs/swagger';

@Exclude()
export class CreateTagsDto {
  @Expose()
  @ApiProperty()
  readonly name: string;
}
