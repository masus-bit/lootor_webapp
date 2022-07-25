import { Exclude, Expose } from 'class-transformer';

@Exclude()
export class SignUpDto {
  @Expose({ name: 'login' })
  readonly login: string;

  @Expose({ name: 'userName' })
  readonly user_name: string;

  @Expose()
  readonly password: string;

  @Expose()
  readonly email: string;
}
