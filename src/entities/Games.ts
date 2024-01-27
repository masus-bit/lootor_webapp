import {
  Column,
  Index,
  JoinColumn,
  JoinTable,
  ManyToMany,
  ManyToOne,
  PrimaryColumn,
} from 'typeorm';
import { Publishers } from './Publishers';
import { Genres } from './Genres';
import { Platforms } from './Platforms';
import { Developers } from './Developers';

// @Entity()
export class Games {
  @PrimaryColumn()
  id: number;

  @Index()
  @Column({
    nullable: true,
  })
  game_title: string;

  @Index()
  @Column({
    nullable: true,
  })
  players: string;

  @Index()
  @Column({ nullable: true })
  release_date: Date;

  @Index()
  @Column({
    nullable: true,
  })
  overview: string;

  @Index()
  @Column({ nullable: true })
  rating: string;

  @Index()
  @Column({
    nullable: true,
  })
  coop: string;

  @ManyToMany(() => Publishers, (pub) => pub.games)
  @JoinTable({ name: 'games_publishers' })
  pubs: Publishers[];

  @ManyToMany(() => Genres, (genre) => genre.games)
  @JoinTable({ name: 'games_genres' })
  genres: Genres[];

  @ManyToOne(() => Platforms, (platform) => platform.id)
  @JoinColumn({ name: 'platform' })
  platform: number;

  @ManyToMany(() => Developers, (devs) => devs.games)
  @JoinTable({ name: 'games_developers' })
  devs: Developers[];
}
