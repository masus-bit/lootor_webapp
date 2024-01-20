import { Exclude, Expose } from 'class-transformer';
import { CollectionItem } from '../../entities/CollectionItem';
import { ApiProperty } from '@nestjs/swagger';
import { ApiModelProperty } from '@nestjs/swagger/dist/decorators/api-model-property.decorator';

@Exclude()
export class CollectionItemCreateDto {
  @Expose()
  @ApiProperty()
  readonly name: string;

  @Expose()
  @ApiProperty()
  readonly description: string;

  @Expose({ name: 'gameId' })
  @ApiProperty({ name: 'gameId' })
  readonly game_id?: string;

  @Expose()
  @ApiProperty()
  readonly developer: string;

  @Expose()
  @ApiProperty()
  readonly publisher: string;

  @Expose()
  @ApiProperty()
  readonly platform: string;

  @Expose({ name: 'itemPicture' })
  @ApiProperty({ name: 'itemPicture' })
  readonly item_picture: string;

  @Expose({ name: 'itemPhotos' })
  @ApiProperty({ name: 'itemPhotos' })
  readonly item_photos: string[];

  @Expose({ name: 'purchaseDate' })
  @ApiProperty({ name: 'purchaseDate' })
  readonly purchase_date: Date;

  @Expose()
  @ApiProperty()
  readonly price: number;

  @Expose()
  @ApiProperty()
  readonly sealed: boolean;

  @Expose({ name: 'limitedEdition' })
  @ApiProperty({ name: 'limitedEdition' })
  readonly limited_edition: boolean;

  @Expose({ name: 'copyNumber' })
  @ApiProperty({ name: 'copyNumber' })
  readonly copy_number?: string;

  @Expose({ name: 'collection' })
  @ApiProperty({ name: 'collection' })
  readonly collection: string;
}

export class ReturnedCollectionItemDto {
  @ApiProperty()
  readonly id?: string;
  @ApiProperty()
  readonly name: string;
  @ApiProperty()
  readonly description: string;
  @ApiProperty()
  readonly gameId?: string;
  @ApiProperty()
  readonly developer: string;
  @ApiProperty()
  readonly publisher: string;
  @ApiProperty()
  readonly platform: string;
  @ApiProperty()
  readonly itemPicture: string;
  @ApiProperty()
  readonly itemPhotos: string[];
  @ApiProperty()
  readonly purchaseDate: Date;
  @ApiProperty({ description: 'Цена' })
  readonly price: number;
  @ApiProperty()
  readonly sealed: boolean;
  @ApiProperty()
  readonly limitedEdition: boolean;
  @ApiProperty()
  readonly copyNumber?: string;
  @ApiProperty()
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
  @ApiModelProperty({ type: ReturnedCollectionItemDto })
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
