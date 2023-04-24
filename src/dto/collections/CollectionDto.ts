import { Collection } from '../../entities/Collection';
import { ReturnUser } from './CreateCollectionDto';
import { ReturnedCollectionItemDto } from '../collectionItem/CollectionItemCreateDto';
import { ApiModelProperty } from '@nestjs/swagger/dist/decorators/api-model-property.decorator';
import { ApiProperty } from '@nestjs/swagger';

export class CollectionDto {
  @ApiProperty()
  readonly id: string;
  @ApiProperty()
  readonly name: string;
  @ApiProperty()
  readonly description: string;
  @ApiModelProperty({ type: ReturnUser })
  readonly user: string;
  @ApiProperty()
  readonly isPrivate: boolean;
  @ApiProperty()
  readonly bannerUrl: string;
  @ApiModelProperty({ type: ReturnedCollectionItemDto, isArray: true })
  readonly collectionItems: ReturnedCollectionItemDto[];
  @ApiProperty()
  readonly collectionItemsCount: number;
  @ApiProperty()
  readonly canLike: boolean;
  @ApiProperty()
  readonly likes: number;
  @ApiProperty()
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
  @ApiModelProperty({ type: CollectionDto })
  readonly data: Readonly<CollectionDto>;

  constructor(data: Readonly<CollectionDto>) {
    this.data = data;
  }
}
