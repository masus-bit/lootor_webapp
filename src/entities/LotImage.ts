import {
  Column,
  Entity,
  Index,
  JoinColumn,
  ManyToOne,
  PrimaryGeneratedColumn,
} from 'typeorm';
import { Lot } from './Lot';
import { DraftLot } from './DraftLot';

@Entity()
export class LotImage {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @ManyToOne(() => Lot, (lot) => lot.lot_id)
  @JoinColumn({ name: 'lot_id' })
  lot_id: string;

  @ManyToOne(() => DraftLot, (draftLot) => draftLot.draft_lot_id)
  @JoinColumn({ name: 'draft_lot_id' })
  draft_lot_id: string;

  @Index()
  @Column({
    nullable: false,
  })
  link: string;
}
