import { Exclude, Expose } from 'class-transformer';

@Exclude()
export class TokenPayloadDto {
  @Expose({ name: 'login' })
  login?: string;

  @Expose()
  iat: number;

  @Expose()
  exp: number;

  @Expose({ name: 'user_name' })
  userName?: string;

  @Expose()
  email?: string;

  @Expose()
  likes?: number;

  @Expose()
  dislikes?: number;

  @Expose()
  created: string;

  @Expose()
  subscriptions?: string[];

  @Expose({ name: 'avatar_url' })
  avatarUrl: string;
}

@Exclude()
export class AuthDto {
  @Expose({ name: 'login' })
  login?: string;

  @Expose()
  iat: number;

  @Expose()
  exp: number;

  @Expose({ name: 'user_name' })
  userName?: string;

  @Expose()
  email?: string;

  @Expose()
  user_name?: string;

  @Expose()
  likes?: number;

  @Expose()
  dislikes?: number;

  @Expose()
  created: string;
}
