import { EntityModel } from '../../entities/EntityModel';

export class GetEntitiesDto {
  readonly data: ReadonlyArray<EntityModel>;

  constructor(data: ReadonlyArray<EntityModel>) {
    this.data = data;
  }
}

export class GetEntityDto {
  readonly data: Readonly<EntityModel>;

  constructor(data: Readonly<EntityModel>) {
    this.data = data;
  }
}
