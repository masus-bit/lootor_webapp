import { Exclude, Expose } from 'class-transformer';
import { Collection } from '../../entities/Collection';
import { User } from '../../entities/User';

@Exclude()
export class CreateCollectionDto {
  @Expose()
  readonly name: string;

  @Expose()
  readonly description: string;

  @Expose({ name: 'isPrivate' })
  readonly is_private: boolean;

  @Expose({ name: 'userLogin' })
  readonly user: string;

  @Expose({ name: 'bannerUrl' })
  readonly banner_url: string;
}

export class ReturnUser {
  readonly login: string;
  readonly userName: string;
  readonly email: string;
  readonly subscriptions: string[];

  constructor(user: Readonly<User>) {
    this.login = user.login;
    this.userName = user.user_name;
    this.email = user.email;
    this.subscriptions = user.subscriptions;
  }
}

export class ReturnedCollectionDto {
  readonly id: string;
  readonly name: string;
  readonly description: string;
  readonly user: string;
  readonly isPrivate: boolean;
  readonly bannerUrl: string;
  readonly likes: number;

  constructor(collection: Readonly<Collection>) {
    this.id = collection.id;
    this.isPrivate = collection.is_private;
    this.name = collection.name;
    this.description = collection.description;
    this.bannerUrl = collection.banner_url;
    this.likes = collection?.likes?.length;
    // @ts-ignore
    this.user = new ReturnUser(collection.user);
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
    this.data = data;
  }
}
