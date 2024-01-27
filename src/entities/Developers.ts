import { Column, Entity, Index, ManyToMany, PrimaryColumn } from 'typeorm';
import { Games } from './Games';

@Entity()
export class Developers {
  @PrimaryColumn()
  id: number;

  @Index()
  @Column({
    nullable: true,
  })
  name: string;

  @ManyToMany(() => Games, (game) => game.devs)
  games: Games[];
}
