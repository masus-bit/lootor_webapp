import { Column, Entity, Index, OneToMany, PrimaryColumn } from 'typeorm';
import { Lot } from './Lot';
import { DraftLot } from './DraftLot';

@Entity()
export class User {
  @PrimaryColumn()
  login: string;

  @Index()
  @Column({
    nullable: false,
    length: 100,
  })
  user_name: string;

  @Index()
  @Column({
    nullable: true,
    length: 100,
    select: false,
  })
  password: string;

  @Index()
  @Column({
    nullable: true,
    length: 100,
    select: false,
  })
  vk_id: string;

  @Index()
  @Column({
    nullable: true,
    length: 100,
    select: false,
  })
  telegram_id: string;

  @Index()
  @Column({
    nullable: true,
    length: 100,
  })
  email: string;

  @Index()
  @Column({
    nullable: true,
  })
  created: Date;

  @Index()
  @Column({
    nullable: true,
    default: 0,
  })
  likes: number;

  @Index()
  @Column({
    nullable: true,
    default: 0,
  })
  dislikes: number;

  @Index()
  @Column({
    nullable: true,
    length: 999,
  })
  avatar_url: string;

  @Index()
  @Column({
    nullable: true,
    length: 100,
  })
  background_url: string;

  @Index()
  @Column({
    nullable: true,
    length: 100,
  })
  verification_token: string;

  @Index()
  @Column({
    nullable: true,
    default: 0,
  })
  subscribers: number;

  @Index()
  @Column('simple-array', { array: true, nullable: false, default: [] })
  subscriptions: string[];

  @Index()
  @Column('simple-array', { array: true, nullable: false, default: [] })
  collection_subscriptions: string[];

  @OneToMany(() => Lot, (lot) => lot.created_by)
  lots_of: Lot[];

  @OneToMany(() => DraftLot, (draftLot) => draftLot.created_by)
  draft_lots_of: Lot[];
}
