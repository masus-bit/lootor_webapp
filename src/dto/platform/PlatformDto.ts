import { ApiProperty } from '@nestjs/swagger';
import { Platforms } from '../../entities/Platforms';

export class PlatformDto {
  @ApiProperty()
  readonly id?: string;
  @ApiProperty()
  readonly name: string;

  constructor(platform: Readonly<Platforms>) {
    this.id = platform?.id;
    this.name = platform?.name;
  }
}

export class GetPlatformsDto {
  readonly data: ReadonlyArray<Platforms>;

  constructor(data: ReadonlyArray<Platforms>) {
    this.data = data;
  }
}
