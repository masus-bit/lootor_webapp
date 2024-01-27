import { Column, Entity, Index, ManyToMany, PrimaryColumn } from 'typeorm';
import { Games } from './Games';

@Entity()
export class Publishers {
  @PrimaryColumn()
  id: number;

  @Index()
  @Column({
    nullable: true,
  })
  keyword: string;

  @Index()
  @Column({
    nullable: true,
  })
  logo: string;

  @ManyToMany(() => Games, (game) => game.pubs)
  games: Games[];
}
