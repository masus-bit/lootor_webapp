import { ApiProperty } from '@nestjs/swagger';
import { ApiModelProperty } from '@nestjs/swagger/dist/decorators/api-model-property.decorator';

export class SignInDto {
  @ApiProperty()
  email: string;
  @ApiProperty()
  password: string;
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
}

export class SignInResponse {
  @ApiModelProperty({ type: AccessTokenPayload })
  accessToken: AccessTokenPayload;
}
