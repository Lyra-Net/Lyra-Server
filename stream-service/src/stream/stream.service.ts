import { Injectable } from '@nestjs/common';
import { Response } from 'express';
import * as ytdl from 'ytdl-core';

@Injectable()
export class StreamService {
  async streamAudioByVideoId(videoId: string, res: Response): Promise<void> {
    const url = `https://www.youtube.com/watch?v=${videoId}`;

    try {
      const info = await ytdl.getInfo(url);
      const format = ytdl.chooseFormat(info.formats, {
        quality: 'highestaudio',
        filter: 'audioonly',
      });

      res.setHeader('Content-Type', format.mimeType || 'audio/webm');
      res.setHeader('Content-Disposition', 'inline');
      res.setHeader('Cache-Control', 'no-cache');

      ytdl(url, {
        quality: 'highestaudio',
        filter: 'audioonly',
        highWaterMark: 1 << 25,
      }).pipe(res);
    } catch (err) {
      throw new Error('Cannot stream audio: ' + (err as Error).message);
    }
  }
}
