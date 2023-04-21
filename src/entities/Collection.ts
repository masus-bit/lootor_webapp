import {
  Column,
  Entity,
  Index,
  JoinColumn,
  ManyToOne,
  OneToMany,
  PrimaryGeneratedColumn,
} from 'typeorm';
import { User } from './User';
import { CollectionItem } from './CollectionItem';

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

  @OneToMany(
    () => CollectionItem,
    (collectionItem) => collectionItem.collection,
  )
  collectionItems: CollectionItem[];

  @ManyToOne(() => User, (user) => user.login, { onDelete: 'CASCADE' })
  @JoinColumn({ name: 'user' })
  user: string;
}
