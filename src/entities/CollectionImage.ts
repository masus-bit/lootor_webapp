import {
  Column,
  Entity,
  Index,
  JoinColumn,
  ManyToOne,
  PrimaryGeneratedColumn,
} from 'typeorm';
import { CollectionItem } from './CollectionItem';

@Entity()
export class CollectionImage {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @ManyToOne(() => CollectionItem, (collection_item) => collection_item.id)
  @JoinColumn({ name: 'collection_item_id' })
  collection_item_id: string;

  @Index()
  @Column({
    nullable: false,
  })
  link: string;
}
