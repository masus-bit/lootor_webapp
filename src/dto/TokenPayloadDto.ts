import { Exclude, Expose } from 'class-transformer';

@Exclude()
export class TokenPayloadDto {
  @Expose({ name: 'id' })
  userId: string;

  @Expose()
  iat: number;

  @Expose()
  exp: number;

  @Expose({ name: 'user_name' })
  userName: string;
}
