import { ApiProperty } from '@nestjs/swagger';

export class SearchDto {
  @ApiProperty()
  readonly index: string;
  @ApiProperty()
  readonly result: { [key: string]: any };

  constructor(searchResult: Readonly<any>) {
    this.index = searchResult._index;
    this.result = searchResult._source;
  }
}

export class SearchResultDto {
  readonly data: ReadonlyArray<SearchDto>;

  constructor(data: ReadonlyArray<SearchDto>) {
    this.data = data;
  }
}
