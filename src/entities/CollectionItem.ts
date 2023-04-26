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
export class CollectionItem {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Index()
  @Column({
    nullable: true,
  })
  game_id: string;

  @Index()
  @Column({ nullable: false })
  name: string;

  @Index()
  @Column({ nullable: false })
  developer: string;

  @Index()
  @Column({ nullable: false })
  publisher: string;

  @Index()
  @Column({ nullable: false })
  platform: string;

  @Index()
  @Column({ nullable: false })
  item_picture: string;

  @Index()
  @Column('simple-array', { array: true, nullable: true })
  item_photos: string[];

  @Index()
  @Column({ nullable: false })
  purchase_date: Date;

  @Index()
  @Column({ nullable: false })
  price: string;

  @Index()
  @Column()
  sealed: boolean;

  @Index()
  @Column()
  limited_edition: boolean;

  @Index()
  @Column({ nullable: true })
  copy_number: string;

  @Index()
  @Column({ default: false, select: false })
  deleted: boolean;

  @Index()
  @Column({ nullable: true })
  @ManyToOne(() => Collection, (collection) => collection.id)
  @JoinColumn({ name: 'collection' })
  collection: string;

  @Index()
  @Column({
    nullable: true,
  })
  description: string;
}
