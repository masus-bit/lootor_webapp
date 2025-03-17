import {
  Column,
  Entity,
  Index,
  JoinColumn,
  ManyToOne,
  OneToMany,
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
  @Column({ nullable: false })
  purchase_date: Date;

  @Index()
  @Column({ nullable: false, default: 0 })
  purchase_price: number;

  @Index()
  @Column()
  sealed: boolean;

  @Index()
  @Column({ nullable: true })
  copy_number: string;

  @Index()
  @Column({ default: false, select: false })
  deleted: boolean;

  @Index()
  @Column({ nullable: true })
  @ManyToOne(() => Collection, (collection) => collection.id, {
    onDelete: 'CASCADE',
  })
  @JoinColumn({ name: 'collection' })
  collection: string;

  @Index()
  @Column({ nullable: true })
  @ManyToOne(() => User, (user) => user.login, {
    onDelete: 'CASCADE',
  })
  @JoinColumn({ name: 'owner' })
  owner: string;

  @Index()
  @Column({ nullable: true })
  @ManyToOne(() => Platforms, (platform) => platform.id, {
    onDelete: 'CASCADE',
  })
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

  @OneToMany(() => EntityModel, (entity) => entity.collection_item)
  entities: EntityModel[];
}
