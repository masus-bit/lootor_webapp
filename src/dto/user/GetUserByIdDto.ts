import { User } from '../../entities/User';
import { IsNotEmpty } from 'class-validator';
import { ApiProperty } from '@nestjs/swagger';
import { ApiModelProperty } from '@nestjs/swagger/dist/decorators/api-model-property.decorator';

export class QueryUserDto {
  @IsNotEmpty()
  @ApiProperty()
  readonly login: string;
}

export class UserByIdDto {
  @ApiProperty()
  readonly login: string;

  @ApiProperty()
  readonly userName: string;

  @ApiProperty()
  readonly email: string;

  @ApiProperty()
  readonly created: number;

  @ApiProperty()
  readonly avatarUrl: string;

  @ApiProperty()
  readonly likes: number;

  @ApiProperty()
  readonly dislikes: number;

  @ApiProperty()
  readonly canSubscribe?: boolean;

  @ApiProperty()
  readonly subscribers?: number;

  @ApiProperty()
  readonly collectionItemsCount?: number;

  @ApiProperty()
  readonly collectionsCount?: number;

  @ApiProperty()
  readonly totalSum?: number;

  @ApiProperty()
  readonly backgroundUrl: string;

  constructor(
    user: Readonly<User>,
    collectionsItemsCount?: number,
    collectionCount?: number,
    totalSum?: number,
    canSubscribe?: boolean,
  ) {
    this.login = user.login;
    this.userName = user.user_name;
    this.email = user.email;
    this.created = user.created;
    this.avatarUrl = user.avatar_url;
    this.likes = user.likes;
    this.dislikes = user.dislikes;
    this.canSubscribe = canSubscribe;
    this.subscribers = user.subscribers;
    this.collectionItemsCount = collectionsItemsCount;
    this.collectionsCount = collectionCount;
    this.totalSum = totalSum || 0;
    this.backgroundUrl = user.background_url;
  }
}

export class GetUserByIdDto {
  @ApiModelProperty({ type: UserByIdDto })
  readonly data: Readonly<UserByIdDto>;

  constructor(data: Readonly<UserByIdDto>) {
    this.data = data;
  }
}
