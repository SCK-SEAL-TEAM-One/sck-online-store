import {
  Entity,
  Column,
  PrimaryGeneratedColumn,
  CreateDateColumn,
  UpdateDateColumn,
} from 'typeorm';

@Entity('points')
export class Point {
  @PrimaryGeneratedColumn()
  id: number;

  @Column({ name: 'org_id' })
  orgId: number;

  @Column({ name: 'user_id' })
  userId: number;

  @Column()
  amount: number;

  @CreateDateColumn()
  created: Date;

  @UpdateDateColumn()
  updated: Date;
}
