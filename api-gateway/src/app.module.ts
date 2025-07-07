import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { UserRoute } from './routes/user.route';
import { SongRouter } from './routes/song.route';
import { StreamRoute } from './routes/stream.route';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
    }),
  ],
  controllers: [StreamRoute, UserRoute, SongRouter],
})
export class AppModule {}
