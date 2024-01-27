import { Column, Entity, Index, PrimaryColumn } from 'typeorm';

@Entity()
export class Platforms {
  @PrimaryColumn()
  id: number;

  @Index()
  @Column({
    nullable: true,
  })
  name: string;

  @Index()
  @Column({
    nullable: true,
  })
  alias: string;

  @Index()
  @Column({
    nullable: true,
  })
  icon: string;

  @Index()
  @Column({
    nullable: true,
  })
  console: string;

  @Index()
  @Column({
    nullable: true,
  })
  controller: string;

  @Index()
  @Column({
    nullable: true,
  })
  developer: string;

  @Index()
  @Column({
    nullable: true,
  })
  manufacturer: string;

  @Index()
  @Column({
    nullable: true,
  })
  media: string;

  @Index()
  @Column({
    nullable: true,
  })
  cpu: string;

  @Index()
  @Column({
    nullable: true,
  })
  memory: string;

  @Index()
  @Column({
    nullable: true,
  })
  graphics: string;

  @Index()
  @Column({
    nullable: true,
  })
  sound: string;

  @Index()
  @Column({
    nullable: true,
  })
  maxcontrollers: string;

  @Index()
  @Column({
    nullable: true,
  })
  display: string;

  @Index()
  @Column({
    nullable: true,
  })
  overview: string;
}
