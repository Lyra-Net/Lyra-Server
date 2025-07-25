import { BadRequestException, Injectable } from '@nestjs/common';
import { Response } from 'express';
import * as ytdl from '@distube/ytdl-core';
// import { IncomingMessage } from 'http';

// interface StreamResponse extends IncomingMessage {
//   headers: {
//     'content-type'?: string;
//     [key: string]: string | string[] | undefined;
//   };
// }

@Injectable()
export class StreamService {
  async streamAudioByVideoId(videoId: string, res: Response): Promise<void> {
    if (!this.isValidVideoId(videoId)) {
      throw new BadRequestException('Invalid videoId');
    }

    const youtubeUrl = `https://www.youtube.com/watch?v=${videoId}`;
    try {
      await this.streamFromYoutube(youtubeUrl, res);
    } catch (err) {
      console.error('Stream error:', {
        message: err instanceof Error ? err.message : 'Unknown error',
        stack: err instanceof Error ? err.stack : undefined,
        videoId,
        youtubeUrl,
      });
      throw new BadRequestException(
        'Cannot stream audio: ' +
          (err instanceof Error ? err.message : 'Unknown error'),
      );
    }
  }

  private async streamFromYoutube(
    youtubeUrl: string,
    res: Response,
  ): Promise<void> {
    const options = this.getYtdlOptions();
    try {
      console.log('Attempting to stream:', { youtubeUrl, options });

      // Validate URL format
      try {
        new URL(youtubeUrl);
        console.log('new url success');
      } catch (urlErr) {
        throw new Error(
          `Invalid YouTube URL format: ${youtubeUrl} - ${urlErr}`,
        );
      }

      // Fetch video info to debug
      try {
        const info = await ytdl.getInfo(youtubeUrl, options);
        console.log('Video info:', {
          title: info.videoDetails.title,
          formats: info.formats.map((f) => ({
            itag: f.itag,
            mimeType: f.mimeType,
            url: f.url,
          })),
        });
      } catch (infoErr) {
        console.error('Failed to fetch video info:', infoErr);
      }

      // Fetch stream
      const stream = ytdl(youtubeUrl, options);
      res.setHeader('Content-Type', 'audio/webm');
      res.setHeader('Content-Disposition', 'inline');
      res.setHeader('Cache-Control', 'no-cache');
      res.setHeader('Accept-Ranges', 'bytes');

      // Pipe stream to response
      stream.pipe(res);

      // Handle stream errors and completion
      return new Promise((resolve, reject) => {
        stream
          .on('error', (err) => {
            console.error('Stream pipe error:', {
              message: err.message,
              stack: err.stack,
              youtubeUrl,
            });
            reject(err);
          })
          .on('end', () => {
            console.log('Stream completed for:', youtubeUrl);
            resolve();
          });
      });
    } catch (err) {
      console.error('Failed to stream from YouTube:', {
        message: err instanceof Error ? err.message : 'Unknown error',
        stack: err instanceof Error ? err.stack : undefined,
        youtubeUrl,
      });
      throw new Error(
        'Failed to stream from YouTube: ' +
          (err instanceof Error ? err.message : 'Unknown error'),
      );
    }
  }

  private isValidVideoId(videoId: string): boolean {
    return /^[a-zA-Z0-9_-]{11}$/.test(videoId);
  }

  private getYtdlOptions(): ytdl.downloadOptions {
    return {
      quality: 'highestaudio',
      filter: 'audioonly',
      requestOptions: {
        headers: {
          'User-Agent':
            'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
          'Accept-Language': 'en-US,en;q=0.9',
        },
      },
    };
  }
}
