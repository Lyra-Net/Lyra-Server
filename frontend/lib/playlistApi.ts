import api from "@/lib/api";

export interface UpdatePlaylistPayload {
  playlist_id: string;
  playlist_name?: string;
  is_public?: boolean;
}

export const playlistApi = {
  get: (playlist_id: string) =>
    api.post("/playlist/get", { playlist_id }),

  update: (payload: UpdatePlaylistPayload) =>
    api.post("/playlist/update", payload),

  delete: (playlist_id: string) =>
    api.post("/playlist/delete", { playlist_id }),

  removeSong: (playlist_id: string, song_id: string) =>
    api.post("/playlist/remove-song", { playlist_id, song_id }),

  moveSong: (playlist_id: string, song_id: string, new_position: number) =>
    api.post("/playlist/move-song", { playlist_id, song_id, new_position }),
};
