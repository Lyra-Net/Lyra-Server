import { BadRequestException, Injectable } from '@nestjs/common';
import { Response } from 'express';
import * as ytdl from '@distube/ytdl-core';
import { CachedFormat } from './entities/cached-format.entity';
import { Repository } from 'typeorm';
import { InjectRepository } from '@nestjs/typeorm';
import { IncomingMessage } from 'http';

interface StreamResponse extends IncomingMessage {
  headers: {
    'content-type'?: string;
    [key: string]: string | string[] | undefined;
  };
}

@Injectable()
export class StreamService {
  private readonly CACHE_EXPIRY_MS = 6 * 60 * 60 * 1000; // 6h cache

  constructor(
    @InjectRepository(CachedFormat)
    private readonly formatRepo: Repository<CachedFormat>,
  ) {}

  async streamAudioByVideoId(videoId: string, res: Response): Promise<void> {
    if (!this.isValidVideoId(videoId)) {
      throw new BadRequestException('Invalid videoId');
    }

    const { format } = await this.getAudioFormat(videoId);

    try {
      await this.streamFromYoutube(
        `https://www.youtube.com/watch?v=${videoId}`,
        res,
        format.mimeType,
      );
    } catch (err) {
      console.error('Stream error:', err);
      throw new BadRequestException(
        'Cannot stream audio: ' +
          (err instanceof Error ? err.message : 'Unknown error'),
      );
    }
  }

  private async getAudioFormat(
    videoId: string,
  ): Promise<{ format: ytdl.VideoFormat }> {
    const cached = await this.formatRepo.findOne({ where: { videoId } });
    const needsRefresh = !cached || this.isCacheExpired(cached);

    if (!needsRefresh && cached?.format?.url) {
      return { format: cached.format };
    }

    const info = await ytdl.getInfo(
      `https://www.youtube.com/watch?v=${videoId}`,
    );
    const format = ytdl.chooseFormat(info.formats, {
      quality: 'highestaudio',
      filter: 'audioonly',
    });

    if (!format?.url) {
      throw new Error('No suitable audio format found');
    }

    // save cache
    const sanitizedFormat = {
      url: format.url,
      mimeType: format.mimeType,
      qualityLabel: format.qualityLabel,
      bitrate: format.bitrate,
      itag: format.itag,
    };

    await this.formatRepo.save({
      videoId,
      format: sanitizedFormat,
    });

    return { format: sanitizedFormat };
  }

  private async streamFromYoutube(
    youtubeUrl: string,
    res: Response,
    mimeType?: string,
  ): Promise<void> {
    return new Promise((resolve, reject) => {
      const stream = ytdl(youtubeUrl, {
        quality: 'highestaudio',
        filter: 'audioonly',
      })
        .on('response', (ytRes: StreamResponse) => {
          res.setHeader(
            'Content-Type',
            ytRes.headers['content-type'] || mimeType || 'audio/webm',
          );
          res.setHeader('Content-Disposition', 'inline');
          res.setHeader('Cache-Control', 'no-cache');
          res.setHeader('Accept-Ranges', 'bytes');
        })
        .on('error', reject)
        .on('end', resolve);

      stream.pipe(res);
    });
  }

  private isCacheExpired(cached: CachedFormat): boolean {
    return Date.now() - cached.createdAt.getTime() > this.CACHE_EXPIRY_MS;
  }

  private isValidVideoId(videoId: string): boolean {
    return /^[a-zA-Z0-9_-]{11}$/.test(videoId);
  }
}
