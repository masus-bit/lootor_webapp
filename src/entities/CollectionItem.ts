import {
  Column,
  Entity,
  Index,
  JoinColumn,
  ManyToOne,
  OneToMany,
  PrimaryGeneratedColumn,
} from 'typeorm';
import { CollectionImage } from './CollectionImage';
import { Collection } from './Collection';

@Entity()
export class CollectionItem {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Index()
  @Column({
    nullable: false,
  })
  game_id: string;

  @ManyToOne(() => Collection, (collection) => collection.id)
  @JoinColumn({ name: 'collection' })
  collection: string;

  @Index()
  @Column({
    nullable: false,
  })
  description: string;

  @OneToMany(
    () => CollectionImage,
    (collectionImage) => collectionImage.collection_item_id,
  )
  images: CollectionImage[];
}
