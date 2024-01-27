import { ApiProperty } from '@nestjs/swagger';
import { Banners } from '../../entities/Banners';

export class BannersDto {
  @ApiProperty()
  readonly id?: number;
  @ApiProperty()
  readonly type: string;
  @ApiProperty()
  readonly side: string;
  @ApiProperty()
  readonly gameId: number;
  @ApiProperty()
  readonly filename: string;
  @ApiProperty()
  readonly resolution: string;

  constructor(banner: Readonly<Banners>) {
    this.id = banner?.id;
    this.type = banner.type;
    this.side = banner.side;
    this.gameId = banner.games_id;
    this.filename = banner.filename;
    this.resolution = banner.resolution;
  }
}
