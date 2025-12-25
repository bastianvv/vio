export interface MediaFile {
  id: number;
  library_id: number;
  episode_id?: number;
  movie_id?: number;
  path: string;
  container: string;
  video_codec: string;
  audio_codec: string;
  width: number;
  height: number;
  audio_channels: number;
  duration_sec: number;
  is_missing: boolean;
}
