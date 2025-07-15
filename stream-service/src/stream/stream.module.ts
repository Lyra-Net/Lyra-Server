import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { StreamService } from './stream.service';
import { StreamController } from './stream.controller';
import { CachedFormat } from './entities/cached-format.entity';

@Module({
  imports: [TypeOrmModule.forFeature([CachedFormat])],
  controllers: [StreamController],
  providers: [StreamService],
})
export class StreamModule {}
