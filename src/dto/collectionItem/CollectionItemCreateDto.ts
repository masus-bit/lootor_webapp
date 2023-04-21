import { Exclude, Expose } from 'class-transformer';
import { CollectionItem } from '../../entities/CollectionItem';

@Exclude()
export class CollectionItemCreateDto {
  @Expose()
  readonly name: string;

  @Expose()
  readonly description: string;

  @Expose({ name: 'gameId' })
  readonly game_id?: string;

  @Expose()
  readonly developer: string;

  @Expose()
  readonly publisher: string;

  @Expose()
  readonly platform: string;

  @Expose({ name: 'itemPicture' })
  readonly item_picture: string;

  @Expose({ name: 'itemPhotos' })
  readonly item_photos: string[];

  @Expose({ name: 'purchaseDate' })
  readonly purchase_date: Date;

  @Expose()
  readonly price: string;

  @Expose()
  readonly sealed: boolean;

  @Expose({ name: 'limitedEdition' })
  readonly limited_edition: boolean;

  @Expose({ name: 'copyNumber' })
  readonly copy_number?: string;

  @Expose({ name: 'collection' })
  readonly collection: string;
}

export class ReturnedCollectionItemDto {
  readonly id?: string;
  readonly name: string;
  readonly description: string;
  readonly gameId?: string;
  readonly developer: string;
  readonly publisher: string;
  readonly platform: string;
  readonly itemPicture: string;
  readonly itemPhotos: string[];
  readonly purchaseDate: Date;
  readonly price: string;
  readonly sealed: boolean;
  readonly limitedEdition: boolean;
  readonly copyNumber?: string;
  readonly collection: string;

  constructor(collectionItem: Readonly<CollectionItem>) {
    this.id = collectionItem?.id;
    this.name = collectionItem.name;
    this.description = collectionItem.description;
    this.gameId = collectionItem?.game_id;
    this.developer = collectionItem.developer;
    this.publisher = collectionItem.publisher;
    this.platform = collectionItem.platform;
    this.itemPicture = collectionItem.item_picture;
    this.itemPhotos = collectionItem.item_photos;
    this.purchaseDate = collectionItem.purchase_date;
    this.price = collectionItem.price;
    this.sealed = collectionItem.sealed;
    this.limitedEdition = collectionItem.limited_edition;
    this.copyNumber = collectionItem?.copy_number;
    // @ts-ignore
    this.collection = collectionItem.collection.id;
  }
}

export class ReturnCreateCollectionItem {
  readonly data: Readonly<ReturnedCollectionItemDto>;

  constructor(data: Readonly<ReturnedCollectionItemDto>) {
    this.data = data;
  }
}

export class ReturnCollectionItemsDto {
  readonly data: ReadonlyArray<ReturnedCollectionItemDto>;

  constructor(data: ReadonlyArray<ReturnedCollectionItemDto>) {
    this.data = data;
  }
}
