import {
  Column,
  Entity,
  Index,
  JoinTable,
  ManyToMany,
  PrimaryGeneratedColumn,
} from 'typeorm';
import { Collection } from './Collection';

@Entity()
export class Tags {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Index()
  @Column({
    nullable: false,
  })
  name: string;

  @ManyToMany(() => Collection, (collection) => collection.tags)
  @JoinTable()
  collections: Collection[];
}
