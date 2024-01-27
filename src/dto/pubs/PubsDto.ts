import { ApiProperty } from '@nestjs/swagger';
import { Publishers } from '../../entities/Publishers';

export class PubsDto {
  @ApiProperty()
  readonly id?: number;
  @ApiProperty()
  readonly keyword: string;
  @ApiProperty()
  readonly logo: string;

  constructor(pub: Readonly<Publishers>) {
    this.id = pub?.id;
    this.keyword = pub.keyword;
    this.logo = pub.logo;
  }
}
