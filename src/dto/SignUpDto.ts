import { Exclude, Expose } from 'class-transformer';
import { ApiProperty } from '@nestjs/swagger';

@Exclude()
export class SignUpDto {
  @Expose({ name: 'login' })
  @ApiProperty({ name: 'login' })
  readonly login: string;

  @Expose({ name: 'userName' })
  @ApiProperty({ name: 'userName' })
  readonly user_name: string;

  @Expose()
  @ApiProperty()
  readonly password: string;

  @Expose()
  @ApiProperty()
  readonly email: string;
}

@Exclude()
export class SignUpDtoResponse {
  @Expose({ name: 'login' })
  @ApiProperty({ name: 'login' })
  readonly login: string;

  @Expose({ name: 'userName' })
  @ApiProperty({ name: 'userName' })
  readonly user_name: string;

  @Expose()
  @ApiProperty()
  readonly email: string;

  @Expose()
  @ApiProperty()
  readonly subscriptions: string[] | null;
}
