import { User } from '../../entities/User';
import { IsNotEmpty } from 'class-validator';

export class QueryUserDto {
  @IsNotEmpty()
  readonly login: string;
}

export class UserByIdDto {
  readonly login: string;
  readonly userName: string;
  readonly email: string;
  readonly created: Date;
  readonly avatarUrl: string;
  readonly likes: number;
  readonly dislikes: number;

  constructor(user: Readonly<User>) {
    this.login = user.login;
    this.userName = user.user_name;
    this.email = user.email;
    this.created = user.created;
    this.avatarUrl = user.avatar_url;
    this.likes = user.likes;
    this.dislikes = user.dislikes;
  }
}

export class GetUserByIdDto {
  readonly data: Readonly<UserByIdDto>;

  constructor(data: Readonly<UserByIdDto>) {
    this.data = data;
  }
}
