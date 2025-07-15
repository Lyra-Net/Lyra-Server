import {
  Entity,
  Column,
  PrimaryGeneratedColumn,
  CreateDateColumn,
  Index,
} from 'typeorm';

@Entity()
export class CachedFormat {
  @PrimaryGeneratedColumn()
  id: number;

  @Index({ unique: true })
  @Column()
  videoId: string;

  @Column({ type: 'jsonb' })
  format: {
    url: string;
    mimeType?: string;
    qualityLabel?: string;
    bitrate?: number;
    itag?: number;
  };

  @CreateDateColumn()
  createdAt: Date;
}
