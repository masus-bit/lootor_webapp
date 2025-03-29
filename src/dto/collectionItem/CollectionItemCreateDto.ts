import { Exclude, Expose } from 'class-transformer';
import { CollectionItem } from '../../entities/CollectionItem';
import { ApiProperty } from '@nestjs/swagger';
import { EntityModel } from '../../entities/EntityModel';
import { Collection } from '../../entities/Collection';

@Exclude()
export class CollectionItemCreateDto {
  @Expose()
  @ApiProperty()
  readonly name: string;

  @Expose()
  @ApiProperty()
  readonly description: string;

  @Expose()
  @ApiProperty()
  readonly platform: string;

  @Expose({ name: 'images' })
  @ApiProperty({ name: 'images' })
  readonly images: string[];

  @Expose({ name: 'purchaseDate' })
  @ApiProperty({ name: 'purchaseDate' })
  readonly purchase_date: Date;

  @Expose({ name: 'purchasePrice' })
  @ApiProperty({ name: 'purchasePrice' })
  readonly purchase_price: number;

  @Expose()
  @ApiProperty()
  readonly sealed: boolean;

  @Expose({ name: 'shippingCost' })
  @ApiProperty({ name: 'shippingCost' })
  readonly shipping_cost: number;

  @Expose({ name: 'copyNumber' })
  @ApiProperty({ name: 'copyNumber' })
  readonly copy_number?: string;

  @Expose({ name: 'collection' })
  @ApiProperty({ name: 'collection' })
  collection?: Collection;

  @Expose()
  @ApiProperty()
  readonly entities: { id: string; name: string }[];
}

export class ReturnedCollectionItemDto {
  @ApiProperty()
  readonly id?: string;
  @ApiProperty()
  readonly name: string;
  @ApiProperty()
  readonly description: string;
  @ApiProperty()
  readonly shipping_cost: number;
  @ApiProperty()
  readonly platform: string;
  @ApiProperty()
  readonly images: string[];
  @ApiProperty()
  readonly purchaseDate: Date;
  @ApiProperty({ description: 'Цена' })
  readonly purchasePrice: number;
  @ApiProperty()
  readonly sealed: boolean;
  @ApiProperty()
  readonly copyNumber?: string;
  @ApiProperty()
  readonly collection: string;
  @ApiProperty({ type: EntityModel, isArray: true })
  readonly entities: EntityModel[];

  constructor(collectionItem: Readonly<CollectionItem>) {
    this.id = collectionItem?.id;
    this.name = collectionItem.name;
    this.description = collectionItem.description;
    // @ts-ignore
    this.platform = collectionItem.platform.id;
    this.images = collectionItem.images;
    this.purchaseDate = collectionItem.purchase_date;
    this.purchasePrice = collectionItem.purchase_price;
    this.sealed = collectionItem.sealed;
    this.copyNumber = collectionItem?.copy_number;
    this.collection = collectionItem?.collections?.[0].id || '';
    this.shipping_cost = collectionItem.shipping_cost;
    this.entities = collectionItem.entities;
  }
}

export class ReturnCreateCollectionItem {
  @ApiProperty({ type: ReturnedCollectionItemDto })
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
