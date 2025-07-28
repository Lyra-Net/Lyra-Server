import { BadRequestException, Injectable, Logger } from '@nestjs/common';
import { Response } from 'express';
import * as ytdl from '@distube/ytdl-core';
import { PassThrough, Readable as NodeJSReadableStream } from 'stream';

@Injectable()
export class StreamService {
  private readonly logger = new Logger(StreamService.name);

  async streamAudioByVideoId(
    videoId: string,
    res: Response,
    quality: string = 'highestaudio',
  ): Promise<void> {
    if (!this.isValidVideoId(videoId)) {
      throw new BadRequestException('Invalid videoId');
    }

    const youtubeUrl = `https://www.youtube.com/watch?v=${videoId}`;
    try {
      await this.streamFromYoutube(youtubeUrl, res, { quality });
    } catch (err) {
      this.logger.error('Stream error:', {
        message: err instanceof Error ? err.message : 'Unknown error',
        videoId,
        youtubeUrl,
      });
      throw new BadRequestException('Cannot stream audio');
    }
  }

  private async streamFromYoutube(
    youtubeUrl: string,
    res: Response,
    options: { quality: string },
  ): Promise<void> {
    const ytdlOptions = this.getYtdlOptions(options.quality);
    let stream: NodeJSReadableStream | null = null;
    let bufferStream: PassThrough | null = new PassThrough();

    try {
      new URL(youtubeUrl);
      const info = await ytdl.getInfo(youtubeUrl, ytdlOptions);
      if (!info.formats.some((format) => format.hasAudio && !format.hasVideo)) {
        throw new BadRequestException('Video has no streamable audio');
      }
      const audioFormat: ytdl.videoFormat = ytdl.chooseFormat(
        info.formats,
        ytdlOptions,
      );

      const contentLength =
        audioFormat.contentLength &&
        typeof audioFormat.contentLength === 'string'
          ? parseInt(audioFormat.contentLength, 10)
          : null;

      const rangeHeader = res.req.headers.range;
      let start = 0;
      let end = contentLength ? contentLength - 1 : undefined;

      if (rangeHeader && contentLength) {
        const parts = rangeHeader.replace(/bytes=/, '').split('-');
        start = parseInt(parts[0], 10);
        end = parts[1] ? parseInt(parts[1], 10) : contentLength - 1;
        const chunkSize = end - start + 1;

        res.status(206);
        res.setHeader(
          'Content-Range',
          `bytes ${start}-${end}/${contentLength}`,
        );
        res.setHeader('Content-Length', chunkSize);
      } else if (contentLength) {
        res.setHeader('Content-Length', contentLength);
      }

      this.logger.log(
        `Streaming with Content-Type: ${audioFormat.mimeType || 'audio/webm'}`,
      );
      res.setHeader('Content-Type', audioFormat.mimeType || 'audio/webm');
      res.setHeader('Accept-Ranges', 'bytes');
      res.setHeader('Cache-Control', 'no-cache');

      try {
        stream = await this.getStreamWithRetry(
          youtubeUrl,
          { ...ytdlOptions, range: { start, end } },
          3,
        );
        this.logger.log(
          `Stream initiated for ${youtubeUrl} with range ${start}-${end}, type: ${stream.constructor.name}`,
        );
      } catch (err) {
        this.logger.error(
          `Failed to initiate stream for ${youtubeUrl}: ${err instanceof Error ? err.message : 'Unknown error'}`,
        );
        throw new BadRequestException('Failed to initiate stream');
      }

      const timeout = setTimeout(() => {
        this.logger.error(`Stream timed out for ${youtubeUrl}`);
        stream?.destroy(new Error('Stream timeout'));
        bufferStream?.destroy();
        res.status(504).end();
      }, 60000); // 60-second timeout

      if (bufferStream) {
        bufferStream.pipe(res); // Pipe bufferStream to response
        if (stream) {
          stream.pipe(bufferStream); // Pipe ytdl stream to bufferStream
          stream.on('error', (err) => {
            clearTimeout(timeout);
            this.logger.error(`Stream error for ${youtubeUrl}`, err);
            throw new Error('Streaming failed');
          });
          stream.on('end', () => {
            clearTimeout(timeout);
            this.logger.log(`Stream completed for ${youtubeUrl}`);
            bufferStream?.end();
          });
        }
      }

      return new Promise(() => {
        res.on('close', () => {
          clearTimeout(timeout);
          if (stream) {
            stream.removeAllListeners();
            stream.destroy();
          }
          if (bufferStream) {
            bufferStream.removeAllListeners();
            bufferStream.destroy();
            bufferStream = null;
          }
        });
      });
    } catch (err) {
      this.logger.error(`Failed to stream from YouTube: ${youtubeUrl}`, err);
      throw new Error('Failed to stream audio');
    } finally {
      if (stream) {
        stream.removeAllListeners();
        stream.destroy();
      }
      if (bufferStream) {
        bufferStream.removeAllListeners();
        bufferStream.destroy();
        bufferStream = null;
      }
    }
  }

  private async getStreamWithRetry(
    url: string,
    options: ytdl.downloadOptions,
    maxRetries = 3,
  ): Promise<NodeJSReadableStream> {
    for (let attempt = 1; attempt <= maxRetries; attempt++) {
      try {
        return ytdl(url, options);
      } catch (err) {
        this.logger.warn(
          `Attempt ${attempt} failed for ${url}: ${err instanceof Error ? err.message : 'Unknown error'}`,
        );
        if (attempt === maxRetries) throw err;
        await new Promise((resolve) => setTimeout(resolve, 2000 * attempt)); // Exponential backoff
      }
    }
    throw new Error('Max retries exceeded');
  }

  private isValidVideoId(videoId: string): boolean {
    return /^[a-zA-Z0-9_-]{11}$/.test(videoId);
  }

  private getYtdlOptions(quality: string): ytdl.downloadOptions {
    return {
      quality,
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
