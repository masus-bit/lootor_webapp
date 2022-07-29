import {
  Column,
  Entity,
  Index,
  JoinColumn,
  ManyToOne,
  PrimaryGeneratedColumn,
} from 'typeorm';
import { Collection } from './Collection';

@Entity()
export class CollectionImage {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @ManyToOne(() => Collection, (collection) => collection.id)
  @JoinColumn({ name: 'collection_id' })
  collection_id: string;

  @Index()
  @Column({
    nullable: false,
  })
  link: string;
}
