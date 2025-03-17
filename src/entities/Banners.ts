import { Column, Entity, Index, PrimaryColumn } from 'typeorm';

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
