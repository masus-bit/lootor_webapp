import { Tags } from '../../entities/Tags';

export class GetTagsDto {
  readonly data: ReadonlyArray<Tags>;

  constructor(data: ReadonlyArray<Tags>) {
    this.data = data;
  }
}

export class GetTagDto {
  readonly data: Readonly<Tags>;

  constructor(data: Readonly<Tags>) {
    this.data = data;
  }
}
