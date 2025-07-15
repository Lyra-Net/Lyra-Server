declare module '@distube/ytdl-core' {
  import { IncomingMessage } from 'http';

  export interface VideoFormat {
    url: string;
    mimeType?: string;
    qualityLabel?: string;
    bitrate?: number;
    itag?: number;
    contentLength?: string;
    audioBitrate?: number;
    hasAudio?: boolean;
    hasVideo?: boolean;
    container?: string;
  }

  export interface VideoInfo {
    formats: VideoFormat[];
    videoDetails: {
      title: string;
      lengthSeconds: string;
    };
  }

  export function getInfo(url: string): Promise<VideoInfo>;
  export function chooseFormat(
    formats: VideoFormat[],
    options: { quality: string; filter: string },
  ): VideoFormat;
  export default function ytdl(
    url: string,
    options: { quality: string; filter?: string },
  ): IncomingMessage;
}
