import { Module } from '@nestjs/common';
import { ClientsModule, Transport } from '@nestjs/microservices';
import { join } from 'path';
import { UserRoute } from './routes/user.route';
import { StreamRoute } from './routes/stream.route';
import { SongRouter } from './routes/song.route';

@Module({
  imports: [
    ClientsModule.register([
      {
        name: 'AUTH_PACKAGE',
        transport: Transport.GRPC,
        options: {
          package: 'auth',
          protoPath: join(__dirname, '../proto/auth.proto'),
          url: process.env.USER_SERVICE_URL,
        },
      },
    ]),
  ],
  controllers: [UserRoute, StreamRoute, SongRouter],
})
export class AppModule {}
