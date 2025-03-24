import { Exclude, Expose } from 'class-transformer';
import { ApiProperty } from '@nestjs/swagger';

@Exclude()
export class TelegramSignInDto {
  @Expose()
  @ApiProperty()
  readonly auth_date: string;

  @Expose()
  @ApiProperty()
  readonly username: string;

  @Expose()
  @ApiProperty()
  readonly hash: string;

  @Expose()
  @ApiProperty()
  readonly first_name: string;

  @Expose()
  @ApiProperty()
  readonly id: string;

  @Expose()
  @ApiProperty()
  readonly last_name: string;

  @Expose()
  @ApiProperty()
  readonly photo_url: string;
}