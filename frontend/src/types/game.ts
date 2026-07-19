export type GameState = 'idle' | 'playing' | 'won' | 'lost';

export interface Card {
  photo: string;
  emoji?: string;
  id: number;
  pairId: string;
  flipped: boolean;
  matched: boolean;
}
