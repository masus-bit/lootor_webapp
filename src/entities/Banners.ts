import {
  Column,
  Entity,
  Index,
  JoinColumn,
  ManyToOne,
  PrimaryColumn,
} from 'typeorm';
import { Games } from './Games';

@Entity()
export class Banners {
  @PrimaryColumn()
  id: number;

  @Index()
  @Column({
    nullable: true,
  })
  type: string;

  @Index()
  @Column({
    nullable: true,
  })
  side: string;

  @ManyToOne(() => Games, (game) => game.id, { onDelete: 'CASCADE' })
  @JoinColumn({ name: 'games_id' })
  games_id: number;

  @Index()
  @Column({
    nullable: true,
  })
  filename: string;

  @Index()
  @Column({
    nullable: true,
  })
  resolution: string;
}
