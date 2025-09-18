
export type Playlist = {
  playlist_id: string;
  playlist_name: string;
  owner_id: string;
  is_public: boolean;
  songs: Song[];
}

export type Playlists = Playlist[];

export type Song = {
  song_id: string;
  title: string;
  categories: string[];
  artists: Artist[];
}

export type Artist = {
    id:   number;
    name: string;
}