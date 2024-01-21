import {
  Column,
  Entity,
  Index,
  JoinColumn,
  ManyToMany,
  ManyToOne,
  OneToMany,
  PrimaryGeneratedColumn,
} from 'typeorm';
import { User } from './User';
import { CollectionItem } from './CollectionItem';
import { Tags } from './Tags';

@Entity()
export class Collection {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Index()
  @Column({
    nullable: false,
  })
  name: string;

  @Index()
  @Column({
    nullable: true,
  })
  description: string;

  @Index()
  @Column({ nullable: true, type: 'bigint' })
  created: number;

  @Index()
  @Column()
  is_private: boolean;

  @Index()
  @Column({
    nullable: true,
  })
  banner_url: string;

  @Index()
  @Column('simple-array', { array: true, nullable: true })
  likes: string[];

  @Index()
  @Column({ default: false, select: false })
  deleted: boolean;

  @Index()
  @Column({
    nullable: true,
  })
  transliteration: string;

  @OneToMany(
    () => CollectionItem,
    (collectionItem) => collectionItem.collection,
  )
  collectionItems: CollectionItem[];

  @Index()
  @Column({
    nullable: false,
    default: 0,
  })
  subscribers_count: number;

  @Index()
  @Column({
    nullable: false,
    default: 0,
  })
  totalPrice: number;

  @Index()
  @Column({
    nullable: true,
  })
  share_string: string;

  @ManyToMany(() => Tags, (tag) => tag.collections)
  tags: Tags[];

  @ManyToOne(() => User, (user) => user.login, { onDelete: 'CASCADE' })
  @JoinColumn({ name: 'user' })
  user: string;
}
