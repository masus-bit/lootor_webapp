import { ApiProperty } from '@nestjs/swagger';

export class RatingDto {
  @ApiProperty()
  readonly isLike: boolean;
}
