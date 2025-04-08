import { Expose } from 'class-transformer';
import { ApiProperty } from '@nestjs/swagger';

export class FeedbackAddDto {
  @Expose()
  @ApiProperty()
  readonly title: string;

  @Expose()
  @ApiProperty()
  readonly description: string;
}
