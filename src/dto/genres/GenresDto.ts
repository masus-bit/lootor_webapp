import { ApiProperty } from '@nestjs/swagger';
import { Genres } from '../../entities/Genres';

export class GenresDto {
  @ApiProperty()
  readonly id?: number;
  @ApiProperty()
  readonly genre: string;

  constructor(genre: Readonly<Genres>) {
    this.id = genre?.id;
    this.genre = genre.genre;
  }
}
