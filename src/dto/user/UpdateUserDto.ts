import { Exclude, Expose } from 'class-transformer';

@Exclude()
export class UpdateUserDto {
  @Expose({ name: 'userName' })
  readonly user_name: string;

  @Expose()
  readonly email: string;

  @Expose({ name: 'avatarUrl' })
  readonly avatar_url: string;
}
