export interface GuestbookEntry {
  id: number;
  nickname: string;
  message: string;
  secret: boolean;
  created_at: string;
}

export interface GameScore {
  id: number;
  nickname: string;
  time_ms: number;
  created_at: string;
}

export interface GamePhotoEntry {
  id: number;
  thumbnail_url: string;
  detail_url: string;
}

export interface GameEmojiEntry {
  id: string;
  emoji: string;
}

export interface GamePhotosResponsePhoto {
  type: 'photo';
  photos: GamePhotoEntry[];
  game_token?: string;
}

export interface GamePhotosResponseEmoji {
  type: 'emoji';
  photos: GameEmojiEntry[];
  game_token?: string;
}

export type GamePhotosResponse =
  | GamePhotosResponsePhoto
  | GamePhotosResponseEmoji;

export interface GameRankingsResponse {
  rankings: GameScore[];
  has_played: boolean;
  game_reset_at: string;
}

export interface HallOfFameEntry {
  id: number;
  nickname: string;
  created_at: string;
}
