import {
  Column,
  Entity,
  Index,
  JoinColumn,
  ManyToMany,
  ManyToOne,
  PrimaryGeneratedColumn,
} from 'typeorm';
import { Collection } from './Collection';
import { User } from './User';
import { Platforms } from './Platforms';
import { EntityModel } from './EntityModel';

@Entity()
export class CollectionItem {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Index()
  @Column({ nullable: false })
  name: string;

  @Index()
  @Column('simple-array', { array: true, nullable: true })
  images: string[];

  @Index()
  @Column({ nullable: true })
  purchase_date: Date;

  @Index()
  @Column({ nullable: false, default: 0 })
  purchase_price: number;

  @Index()
  @Column({ nullable: false, default: false })
  sealed: boolean;

  @Index()
  @Column({ nullable: true })
  edition: string;

  @Index()
  @Column('simple-array', { array: true, nullable: true })
  copy_number: string[];

  @Index()
  @Column({ default: false, select: false })
  deleted: boolean;

  @ManyToMany(() => Collection, (collection) => collection.collectionItems, {
    onDelete: 'CASCADE',
  })
  collections: Collection[];

  @Index()
  @ManyToOne(() => User, (user) => user.login, {
    onDelete: 'CASCADE',
  })
  @JoinColumn({ name: 'owner' })
  owner: string;

  @ManyToOne(() => Platforms, (platform) => platform.id, {
    onDelete: 'CASCADE',
  })
  @JoinColumn({ name: 'platform' })
  platform: string;

  @Index()
  @Column({
    nullable: true,
  })
  description: string;

  @Index()
  @Column({
    nullable: true,
  })
  shipping_cost: number;

  @Index()
  @Column({
    nullable: true,
  })
  transliteration: string;

  @ManyToMany(() => EntityModel, (entity) => entity.collection_item)
  entities: EntityModel[];
}
