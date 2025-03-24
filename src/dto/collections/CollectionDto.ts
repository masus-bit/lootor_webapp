import { Collection } from '../../entities/Collection';
import { ReturnUser } from './CreateCollectionDto';
import { ReturnedCollectionItemDto } from '../collectionItem/CollectionItemCreateDto';
import { ApiProperty } from '@nestjs/swagger';
import { Tags } from '../../entities/Tags';

export class CollectionDto {
  @ApiProperty()
  readonly id: string;
  @ApiProperty()
  readonly name: string;
  @ApiProperty()
  readonly description: string;
  @ApiProperty({ type: ReturnUser })
  readonly user: string;
  @ApiProperty()
  readonly isPrivate: boolean;
  @ApiProperty()
  readonly bannerUrl: string;
  @ApiProperty({ type: ReturnedCollectionItemDto, isArray: true })
  readonly collectionItems: ReturnedCollectionItemDto[];
  @ApiProperty()
  readonly collectionItemsCount: number;
  @ApiProperty()
  readonly canLike: boolean;
  @ApiProperty()
  readonly likes: number;
  @ApiProperty()
  readonly transliteration: string;
  @ApiProperty({ type: Tags, isArray: true })
  readonly tags: Tags[];
  @ApiProperty()
  readonly created: number;
  @ApiProperty({ name: 'subscribers_count' })
  readonly subscribersCount: number;
  @ApiProperty()
  readonly totalPrice: number;
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

    // @ts-ignore
    this.user = new ReturnUser(collection.user);
    this.collectionItems =
      collection.collectionItems?.map(
        (item) => new ReturnedCollectionItemDto(item),
      ) || [];
    this.collectionItemsCount = collection.collectionItems?.length || 0;
    // @ts-ignore
    this.canLike = !collection.likes?.includes(collection.user.login);
    this.likes = collection?.likes?.length || 0;
    this.transliteration = collection.transliteration;
    this.tags = collection.tags;
    this.created = collection.created;
    this.subscribersCount = collection.subscribers_count;
    this.totalPrice = collection.totalPrice;
    this.canSubscribe = canSubscribe;
    this.shareString = collection.share_string;
  }
}

export class GetOneCollectionDto {
  @ApiProperty({ type: CollectionDto })
  readonly data: Readonly<CollectionDto>;

  constructor(data: Readonly<CollectionDto>) {
    this.data = data;
  }
}
