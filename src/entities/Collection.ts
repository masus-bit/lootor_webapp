import {
  Column,
  Entity,
  Index,
  JoinColumn,
  ManyToOne,
  PrimaryGeneratedColumn,
} from 'typeorm';
import { User } from './User';

@Entity()
export class Collection {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Index()
  @Column({
    nullable: false,
  })
  name: string;

  @Index()
  @Column({
    nullable: true,
  })
  description: string;

  @Index()
  @Column()
  is_private: boolean;

  @Index()
  @Column({
    nullable: true,
  })
  banner_url: string;

  @ManyToOne(() => User, (user) => user.login, { onDelete: 'CASCADE' })
  @JoinColumn({ name: 'user' })
  user: string;
}
