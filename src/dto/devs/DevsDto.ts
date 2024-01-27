import { ApiProperty } from '@nestjs/swagger';
import { Developers } from '../../entities/Developers';

export class DevsDto {
  @ApiProperty()
  readonly id?: number;
  @ApiProperty()
  readonly name: string;

  constructor(dev: Readonly<Developers>) {
    this.id = dev?.id;
    this.name = dev.name;
  }
}
