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
export class EntityModel {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Index()
  @Column({ nullable: false })
  name: string;

  @Index()
  @Column({ nullable: true })
  @ManyToOne(() => CollectionItem, (collectionItem) => collectionItem.id, {
    onDelete: 'CASCADE',
  })
  @JoinColumn({ name: 'collection_item' })
  collection_item: string;
}
