import { Exclude, Expose } from 'class-transformer';

@Exclude()
export class TokenPayloadDto {
  @Expose({ name: 'login' })
  login: string;

  @Expose()
  iat: number;

  @Expose()
  exp: number;

  @Expose({ name: 'user_name' })
  userName: string;

  @Expose()
  email: string;
}
