import { Exclude, Expose } from 'class-transformer';
import { ApiProperty } from '@nestjs/swagger';

@Exclude()
export class RefreshDto {
  @Expose()
  @ApiProperty({ name: 'refreshToken' })
  readonly refreshToken: string;
}
