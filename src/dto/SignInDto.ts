import { ApiProperty } from '@nestjs/swagger';
;

export class SignInDto {
  @ApiProperty()
  email: string;
  @ApiProperty()
  password: string;

  login?: string
}

class AccessTokenPayload {
  @ApiProperty()
  login: string;
  @ApiProperty()
  userName: string;
  @ApiProperty()
  email: string;
  @ApiProperty()
  likes: number;
  @ApiProperty()
  dislikes: number;
  @ApiProperty()
  subscriptions: string[] | null;
  @ApiProperty()
  avatarUrl: string;
  @ApiProperty()
  created: string;
}

export class SignInResponse {
  @ApiProperty({ type: AccessTokenPayload })
  accessToken: AccessTokenPayload;
}
