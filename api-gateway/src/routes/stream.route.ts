import { Controller, Get, Param, Req, Res } from '@nestjs/common';
import { Request, Response } from 'express';
import { forwardRequest } from '../utils/http';

@Controller()
export class StreamRoute {
  @Get('play/:videoId')
  async streamAudio(
    @Param('videoId') videoId: string,
    @Req() req: Request,
    @Res() res: Response,
  ) {
    const rangeHeader = req.headers.range;

    try {
      const streamRes = await forwardRequest(
        `${process.env.STREAM_SERVICE_URL}/stream/${videoId}`,
        'GET',
        null,
        {
          responseType: 'stream',
          headers: {
            Range: rangeHeader || '',
          },
        },
      );

      res.status(streamRes.status);
      for (const [key, value] of Object.entries(streamRes.headers)) {
        if (
          !['transfer-encoding', 'content-length', 'connection'].includes(
            key.toLowerCase(),
          )
        ) {
          res.setHeader(key, value as string);
        }
      }

      streamRes.data.pipe(res);
    } catch (err: any) {
      console.error('Streaming error via gateway:', err.message);

      if (err.response) {
        res
          .status(err.response.status || 502)
          .json({ message: err.response.data?.message || 'Upstream error' });
      } else {
        res.status(500).json({ message: 'Internal server error' });
      }
    }
  }
}
