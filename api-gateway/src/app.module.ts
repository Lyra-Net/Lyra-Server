import { Module } from '@nestjs/common';
import { ClientsModule, Transport } from '@nestjs/microservices';
import { join } from 'path';
import { ConfigModule, ConfigService } from '@nestjs/config';

import { UserRoute } from './routes/user.route';
import { StreamRoute } from './routes/stream.route';
import { SongRouter } from './routes/song.route';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
    }),

    ClientsModule.registerAsync([
      {
        name: 'AUTH_PACKAGE',
        imports: [ConfigModule],
        inject: [ConfigService],
        useFactory: (configService: ConfigService) => ({
          transport: Transport.GRPC,
          options: {
            package: 'auth',
            protoPath: join(__dirname, '../proto/auth.proto'),
            url: configService.get<string>('USER_SERVICE_URL'),
          },
        }),
      },
    ]),
  ],
  controllers: [UserRoute, StreamRoute, SongRouter],
})
export class AppModule {}
