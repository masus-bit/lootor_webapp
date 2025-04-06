import { ApiProperty } from '@nestjs/swagger';

export class ReturnDataDto {
  @ApiProperty()
  readonly data: any;

  constructor(data: any) {
    this.data = data;
  }
}
