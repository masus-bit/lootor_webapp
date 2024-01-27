import { ApiProperty } from '@nestjs/swagger';
import { Platforms } from '../../entities/Platforms';

export class PlatformDto {
  @ApiProperty()
  readonly id?: number;
  @ApiProperty()
  readonly name: string;
  @ApiProperty()
  readonly alias: string;
  @ApiProperty()
  readonly icon: string;
  @ApiProperty()
  readonly console: string;
  @ApiProperty()
  readonly controller: string;
  @ApiProperty()
  readonly developer: string;
  @ApiProperty()
  readonly manufacturer: string;
  @ApiProperty()
  readonly media: string;
  @ApiProperty()
  readonly cpu: string;
  @ApiProperty()
  readonly memory: string;
  @ApiProperty()
  readonly graphics: string;
  @ApiProperty()
  readonly sound: string;
  @ApiProperty()
  readonly maxcontrollers: string;
  @ApiProperty()
  readonly display: string;
  @ApiProperty()
  readonly overview: string;

  constructor(platform: Readonly<Platforms>) {
    this.id = platform?.id;
    this.name = platform.name;
    this.alias = platform.alias;
    this.icon = platform?.icon;
    this.console = platform?.console;
    this.controller = platform?.controller;
    this.developer = platform?.developer;
    this.manufacturer = platform?.manufacturer;
    this.media = platform?.media;
    this.cpu = platform?.cpu;
    this.memory = platform?.memory;
    this.graphics = platform?.graphics;
    this.sound = platform?.sound;
    this.maxcontrollers = platform?.maxcontrollers;
    this.display = platform?.display;
    this.overview = platform?.overview;
  }
}
