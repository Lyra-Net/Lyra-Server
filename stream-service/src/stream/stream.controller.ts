import {
  Controller,
  Get,
  Param,
  Res,
  Header,
  HttpException,
  HttpStatus,
} from '@nestjs/common';
import { StreamService } from './stream.service';
import { Response } from 'express';

@Controller('stream')
export class StreamController {
  constructor(private readonly streamService: StreamService) {}

  @Get(':videoId')
  @Header('Content-Type', 'audio/webm')
  async streamAudio(@Param('videoId') videoId: string, @Res() res: Response) {
    try {
      await this.streamService.streamAudioByVideoId(videoId, res);
    } catch (err) {
      console.error(err);
      throw new HttpException(
        'Unable to stream audio',
        HttpStatus.INTERNAL_SERVER_ERROR,
      );
    }
  }
}
