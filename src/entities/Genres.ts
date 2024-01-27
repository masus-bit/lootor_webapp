import { Column, Entity, Index, ManyToMany, PrimaryColumn } from 'typeorm';
import { Games } from './Games';

@Entity()
export class Genres {
  @PrimaryColumn()
  id: number;

  @Index()
  @Column({
    nullable: true,
  })
  genre: string;

  @ManyToMany(() => Games, (game) => game.pubs)
  games: Games[];
}
