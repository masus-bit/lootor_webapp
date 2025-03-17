import { Column, Entity, Index, PrimaryColumn } from 'typeorm';

@Entity()
export class Genres {
  @PrimaryColumn()
  id: number;

  @Index()
  @Column({
    nullable: true,
  })
  genre: string;
}
