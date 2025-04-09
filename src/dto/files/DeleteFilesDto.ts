import { ApiProperty } from '@nestjs/swagger';
import { Expose } from 'class-transformer';

export class DeleteFilesDto {
  @ApiProperty()
  @Expose()
  keys: string[];
}
