import { Exclude, Expose } from 'class-transformer';
import { Collection } from '../../entities/Collection';
import { User } from '../../entities/User';
import { ApiProperty } from '@nestjs/swagger';
import { ApiModelProperty } from '@nestjs/swagger/dist/decorators/api-model-property.decorator';
import { Tags } from '../../entities/Tags';

@Exclude()
export class CreateCollectionDto {
  @Expose()
  @ApiProperty()
  readonly name: string;

  @Expose()
  @ApiProperty()
  readonly description: string;

  @Expose({ name: 'isPrivate' })
  @ApiProperty({ name: 'isPrivate' })
  readonly is_private: boolean;

  @Expose({ name: 'userLogin' })
  @ApiProperty({ name: 'userLogin' })
  readonly user: string;

  @Expose({ name: 'bannerUrl' })
  @ApiProperty({ name: 'bannerUrl' })
  readonly banner_url: string;

  @Expose()
  @ApiProperty()
  readonly transliteration: string;

  @Expose()
  @ApiProperty()
  readonly tags: { id: string; name: string }[];
}

export class ReturnUser {
  @ApiProperty()
  readonly login: string;
  @ApiProperty()
  readonly userName: string;
  @ApiProperty()
  readonly email: string;
  @ApiProperty()
  readonly subscriptions: string[];
  @ApiProperty()
  readonly avatarUrl: string;

  constructor(user: Readonly<User>) {
    this.login = user.login;
    this.userName = user.user_name;
    this.email = user.email;
    this.subscriptions = user.subscriptions;
    this.avatarUrl = user.avatar_url;
  }
}

export class ReturnedCollectionDto {
  @ApiProperty()
  readonly id: string;
  @ApiProperty()
  readonly name: string;
  @ApiProperty()
  readonly description: string;
  @ApiModelProperty({ type: ReturnUser })
  readonly user?: string;
  @ApiProperty()
  readonly isPrivate: boolean;
  @ApiProperty()
  readonly bannerUrl: string;
  @ApiProperty()
  readonly likes: number;
  @ApiProperty()
  readonly transliteration: string;
  @ApiProperty()
  readonly canLike: boolean;
  @ApiProperty()
  readonly collectionItemsCount: number;
  @ApiModelProperty({ type: Tags, isArray: true })
  readonly tags: Tags[];
  @ApiProperty()
  readonly created: number;
  @ApiProperty({ name: 'subscribers_count' })
  readonly subscribersCount: number;
  @ApiProperty()
  readonly totalPrice: number;
  @ApiProperty()
  readonly shippingTotal: number;
  @ApiProperty()
  readonly canSubscribe: boolean;
  @ApiProperty()
  readonly shareString: string;

  constructor(collection: Readonly<Collection>, canSubscribe?: boolean) {
    this.id = collection.id;
    this.isPrivate = collection.is_private;
    this.name = collection.name;
    this.description = collection.description;
    this.bannerUrl = collection.banner_url;
    this.likes = collection?.likes?.length || 0;
    this.transliteration = collection.transliteration;
    // @ts-ignore
    this.user = collection.user ? new ReturnUser(collection?.user) : null;
    // @ts-ignore
    this.canLike = !collection.likes?.includes(collection.user?.login);
    this.collectionItemsCount = collection.collectionItems?.length || 0;
    this.tags = collection.tags;
    this.created = collection.created;
    this.subscribersCount = collection.subscribers_count;
    this.totalPrice = collection.totalPrice;
    this.shippingTotal = collection.shippingTotal;
    this.canSubscribe = canSubscribe;
    this.shareString = collection.share_string;
  }
}

export class ReturnCreateCollection {
  @ApiModelProperty({ type: ReturnedCollectionDto })
  readonly data: Readonly<ReturnedCollectionDto>;

  constructor(data: Readonly<ReturnedCollectionDto>) {
    this.data = data;
  }
}

export class ReturnCollectionsDto {
  @ApiModelProperty({ type: ReturnedCollectionDto, isArray: true })
  readonly data: ReadonlyArray<ReturnedCollectionDto>;

  constructor(data: ReadonlyArray<ReturnedCollectionDto>) {
    this.data = data;
  }
}
