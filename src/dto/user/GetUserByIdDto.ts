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
  readonly created: Date;

  @ApiProperty()
  readonly avatarUrl: string;

  @ApiProperty()
  readonly likes: number;

  @ApiProperty()
  readonly dislikes: number;

  @ApiProperty()
  readonly canSubscribe?: boolean;

  constructor(user: Readonly<User>, canSubscribe?: boolean) {
    this.login = user.login;
    this.userName = user.user_name;
    this.email = user.email;
    this.created = user.created;
    this.avatarUrl = user.avatar_url;
    this.likes = user.likes;
    this.dislikes = user.dislikes;
    this.canSubscribe = canSubscribe;
  }
}

export class GetUserByIdDto {
  @ApiModelProperty({ type: UserByIdDto })
  readonly data: Readonly<UserByIdDto>;

  constructor(data: Readonly<UserByIdDto>) {
    this.data = data;
  }
}
