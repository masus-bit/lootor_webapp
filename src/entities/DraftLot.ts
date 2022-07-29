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
import { LotImage } from './LotImage';

@Entity()
export class DraftLot {
  @PrimaryGeneratedColumn('uuid')
  draft_lot_id: string;

  @Index()
  @Column({
    nullable: false,
  })
  platform: string;

  @Index()
  @Column({
    nullable: true,
  })
  region: string;

  @Index()
  @Column({
    nullable: false,
  })
  bid: number;

  @Index()
  @Column({
    nullable: false,
  })
  description: string;

  @Index()
  @Column({
    nullable: false,
  })
  bid_step: number;

  @Index()
  @Column({
    nullable: false,
  })
  created_at: Date;

  @Index()
  @Column({
    nullable: false,
  })
  is_sealed: boolean;

  @Index()
  @Column({
    nullable: false,
  })
  condition: number;

  @Index()
  @Column({
    nullable: true,
  })
  @ManyToOne(() => User, (user) => user.login)
  @JoinColumn({ name: 'created_by' })
  created_by: string;

  @OneToMany(() => LotImage, (lotImage) => lotImage.draft_lot_id)
  images: LotImage[];
}
