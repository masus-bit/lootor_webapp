import { Column, Entity, Index, PrimaryColumn } from 'typeorm';

@Entity()
export class User {
  @PrimaryColumn()
  login: string;

  @Index()
  @Column({
    nullable: false,
    length: 100,
  })
  user_name: string;

  @Index()
  @Column({
    nullable: true,
    length: 100,
    select: false
  })
  password: string;

  @Index()
  @Column({
    nullable: true,
    length: 100,
    select: false
  })
  password_encrypted: string;

  @Index()
  @Column({
    nullable: true,
    length: 100,
  })
  email: string;
}
