import { Collection } from '../../entities/Collection';
import { ReturnUser } from './CreateCollectionDto';
import { ReturnedCollectionItemDto } from '../collectionItem/CollectionItemCreateDto';

export class CollectionDto {
  readonly id: string;
  readonly name: string;
  readonly description: string;
  readonly user: string;
  readonly isPrivate: boolean;
  readonly bannerUrl: string;
  readonly collectionItems: ReturnedCollectionItemDto[];
  readonly collectionItemsCount: number;
  readonly canLike: boolean;
  readonly likes: number;
  readonly transliteration: string;

  constructor(collection: Readonly<Collection>) {
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
  }
}

export class GetOneCollectionDto {
  readonly data: Readonly<CollectionDto>;

  constructor(data: Readonly<CollectionDto>) {
    this.data = data;
  }
}
