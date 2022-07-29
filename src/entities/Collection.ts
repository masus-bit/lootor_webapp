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
import { CollectionImage } from './CollectionImage';

@Entity()
export class Collection {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Index()
  @Column({
    nullable: false,
  })
  game_id: string;

  @ManyToOne(() => User, (user) => user.login)
  @JoinColumn({ name: 'user' })
  user: string;

  @Index()
  @Column({
    nullable: false,
  })
  description: string;

  @OneToMany(
    () => CollectionImage,
    (collectionImage) => collectionImage.collection_id,
  )
  images: CollectionImage[]
}
