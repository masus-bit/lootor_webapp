import { ApiProperty } from '@nestjs/swagger';
import { Platforms } from '../../entities/Platforms';

export class PlatformDto {
  @ApiProperty()
  readonly id?: string;
  @ApiProperty()
  readonly name: string;

  constructor(platform: Readonly<Platforms>) {
    this.id = platform?.id;
    this.name = platform.name;
  }
}
