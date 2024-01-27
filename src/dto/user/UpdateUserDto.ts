import { Exclude, Expose } from 'class-transformer';
import { ApiProperty } from '@nestjs/swagger';

@Exclude()
export class UpdateUserDto {
  @Expose({ name: 'userName' })
  @ApiProperty({ name: 'userName' })
  readonly user_name: string;

  @Expose()
  @ApiProperty()
  readonly email: string;

  @Expose({ name: 'avatarUrl' })
  @ApiProperty({ name: 'avatarUrl' })
  readonly avatar_url: string;

  @Expose({ name: 'backgroundUrl' })
  @ApiProperty({ name: 'backgroundUrl' })
  readonly background_url: string;
}
