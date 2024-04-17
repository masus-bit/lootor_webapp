import {
  Column,
  Entity,
  Index,
  JoinColumn,
  ManyToOne,
  PrimaryGeneratedColumn,
} from 'typeorm';
import { User } from './User';
import { CollectionItem } from './CollectionItem';
import { Collection } from './Collection';

@Entity()
export class Event {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Index()
  @Column({
    nullable: false,
  })
  action: string;

  @Index()
  @Column({ nullable: true, type: 'bigint' })
  date: number;

  @Index()
  @Column({ nullable: true })
  event_target: string;

  @Index()
  @Column({ nullable: true })
  target_name: string;

  @ManyToOne(() => User, (user) => user.login, { onDelete: 'CASCADE' })
  @JoinColumn({ name: 'user' })
  user: string;

  @ManyToOne(() => User, (user) => user.login, {
    onDelete: 'SET NULL',
    nullable: true,
  })
  @JoinColumn({ name: 'target_user' })
  target_user: string;

  @ManyToOne(() => CollectionItem, (item) => item.id, {
    onDelete: 'SET NULL',
    nullable: true,
  })
  @JoinColumn({ name: 'target_collection_item' })
  target_collection_item: string;

  @ManyToOne(() => Collection, (collection) => collection.id, {
    onDelete: 'SET NULL',
    nullable: true,
  })
  @JoinColumn({ name: 'target_collection' })
  target_collection: string;
}
