import { ApiProperty } from '@nestjs/swagger';
import { ApiModelProperty } from '@nestjs/swagger/dist/decorators/api-model-property.decorator';
import { ReturnUser } from '../collections/CreateCollectionDto';
import { GenresDto } from '../genres/GenresDto';
import { Games } from '../../entities/Games';
import { PubsDto } from '../pubs/PubsDto';
import { PlatformDto } from '../platform/PlatformDto';
import { DevsDto } from '../devs/DevsDto';
import { BannersDto } from '../banners/BannersDto';

export class GamesDto {
  @ApiProperty()
  readonly id: number;
  @ApiProperty()
  readonly gameTitle: string;
  @ApiProperty()
  readonly players: string;
  @ApiModelProperty({ type: ReturnUser })
  readonly releaseDate: Date;
  @ApiProperty()
  readonly overview: string;
  @ApiProperty()
  readonly rating: string;
  @ApiModelProperty({ type: GenresDto, isArray: true })
  readonly genre: GenresDto;
  @ApiProperty()
  readonly coop: string;
  @ApiProperty()
  readonly publisher: PubsDto;
  @ApiProperty()
  readonly platform: PlatformDto;
  @ApiProperty()
  readonly developer: DevsDto;
  @ApiProperty()
  readonly banners: BannersDto[];

  constructor(game: Readonly<Games>) {
    this.id = game.id;
    this.gameTitle = game.game_title;
    this.players = game.players;
    this.releaseDate = game.release_date;
    this.overview = game.overview;
    this.rating = game.rating;
    this.genre = game.genres.map((item) => new GenresDto(item))[0];
    this.overview = game.overview;
    this.coop = game.coop;
    this.publisher = game.pubs.map((item) => new PubsDto(item))[0];
    // @ts-ignore
    this.platform = new PlatformDto(game.platform);
    this.developer = game.devs.map((item) => new DevsDto(item))[0];
    // @ts-ignore
    this.banners = game.banners.map((item) => new BannersDto(item));
  }
}

export class ReturnArrayGamesDto {
  @ApiModelProperty({ type: GamesDto, isArray: true })
  readonly data: ReadonlyArray<GamesDto>;

  constructor(data: ReadonlyArray<GamesDto>) {
    this.data = data;
  }
}
