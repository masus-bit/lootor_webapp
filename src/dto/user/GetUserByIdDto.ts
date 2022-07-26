import { User } from '../../entities/User';
import { IsNotEmpty } from 'class-validator';

export class QueryUserDto {
  @IsNotEmpty()
  readonly login: string;
}

export class UserByIdDto {
  readonly login: string;
  readonly user_name: string;
  readonly email: string;

  constructor(user: Readonly<User>) {
    this.login = user.login;
    this.user_name = user.user_name;
    this.email = user.email;
  }
}

export class GetUserByIdDto {
  readonly data: Readonly<UserByIdDto>;

  constructor(data: Readonly<UserByIdDto>) {
    this.data = data;
  }
}
