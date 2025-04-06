import {
  Column,
  Entity,
  Index,
  JoinColumn,
  JoinTable,
  ManyToMany,
  ManyToOne,
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
  @Column({ nullable: true })
  created: Date;

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

  @ManyToMany(
    () => CollectionItem,
    (collectionItem) => collectionItem.collections,
  )
  @JoinTable()
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
    nullable: false,
    default: 0,
  })
  shippingTotal: number;

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
