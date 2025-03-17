import {
  Column,
  Entity,
  Index,
  JoinTable,
  ManyToMany,
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

  @ManyToMany(
    () => CollectionItem,
    (collectionItem) => collectionItem.entities,
    {
      onDelete: 'CASCADE',
    },
  )
  @JoinTable()
  collection_item: CollectionItem[];
}
