import { Exclude, Expose } from 'class-transformer';
import { Collection } from '../../entities/Collection';

class Attributes {}

@Exclude()
export class CreateCollectionDto {
  @Expose()
  readonly name: string;

  @Expose()
  readonly description: string;

  @Expose({name: 'isPrivate'})
  readonly is_private: boolean;

  @Expose({name: 'userLogin'})
  readonly user: string;
}

export class ReturnedCollectionDto {
  readonly id: string;
  readonly name: string;
  readonly description: string;
  readonly user: string;
  readonly isPrivate: boolean;

  constructor(collection: Readonly<Collection>) {
    console.log(collection);
    this.id = collection.id;
    this.isPrivate = collection.is_private;
    this.name = collection.name;
    this.description = collection.description;
    this.user = collection.user;
  }
}

export class ReturnCreateCollection {
  readonly data: Readonly<ReturnedCollectionDto>;

  constructor(data: Readonly<ReturnedCollectionDto>) {
    this.data = data;
  }
}

export class ReturnCollectionsDto {
  readonly data: ReadonlyArray<ReturnedCollectionDto>;

  constructor(data: ReadonlyArray<ReturnedCollectionDto>) {
    this.data = data
  }
}