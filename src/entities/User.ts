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
  password_encrypted: string;

  @Index()
  @Column({
    nullable: true,
    length: 100,
  })
  email: string;

  @OneToMany(() => Lot, (lot) => lot.created_by)
  lots_of: Lot[];

  @OneToMany(() => DraftLot, (draftLot) => draftLot.created_by)
  draft_lots_of: Lot[];
}
