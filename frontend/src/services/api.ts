import type {
  GuestbookEntry,
  GameScore,
  GameRankingsResponse,
  GamePhotosResponse,
  PhotoUpload,
  HallOfFameEntry,
} from '../types';
import { i18n } from '../i18n';

const API_BASE = '/api';

export class ApiError extends Error {
  code: string | null;
  status: number;

  constructor(message: string, code: string | null, status: number) {
    super(message);
    this.code = code;
    this.status = status;
  }
}

async function fail(res: Response, fallback: string): Promise<never> {
  const body = await res.json().catch(() => null);
  throw new ApiError(body?.error ?? fallback, body?.code ?? null, res.status);
}

export function apiErrorMessage(e: unknown, fallback?: string): string {
  const { t, te } = i18n.global;
  if (e instanceof ApiError) {
    if (e.code && te(`errors.${e.code}`)) return t(`errors.${e.code}`);
    return fallback ?? e.message;
  }
  return fallback ?? t('errors.network');
}

export interface PaginatedGuestbook {
  items: GuestbookEntry[];
  next_cursor: number | null;
}

export async function getGuestbook(
  cursor?: number,
  limit = 10,
): Promise<PaginatedGuestbook> {
  const params = new URLSearchParams({ limit: String(limit) });
  if (cursor) params.set('cursor', String(cursor));
  const res = await fetch(`${API_BASE}/guestbook?${params}`);
  if (!res.ok) await fail(res, 'Failed to fetch guestbook');
  return res.json();
}

export async function createGuestbookEntry(
  nickname: string,
  message: string,
  password: string,
  secret: boolean = false,
): Promise<GuestbookEntry> {
  const res = await fetch(`${API_BASE}/guestbook`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ nickname, message, password, secret }),
  });
  if (!res.ok) await fail(res, 'Failed to create entry');
  return res.json();
}

export async function verifyGuestbookPassword(
  id: number,
  password: string,
): Promise<void> {
  const res = await fetch(`${API_BASE}/guestbook/${id}/verify`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ password }),
  });
  if (!res.ok) await fail(res, 'Failed to verify password');
}

export async function updateGuestbookEntry(
  id: number,
  message: string,
  password: string,
): Promise<GuestbookEntry> {
  const res = await fetch(`${API_BASE}/guestbook/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ message, password }),
  });
  if (!res.ok) await fail(res, 'Failed to update entry');
  return res.json();
}

export async function deleteGuestbookEntry(
  id: number,
  password: string,
): Promise<void> {
  const res = await fetch(`${API_BASE}/guestbook/${id}`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ password }),
  });
  if (!res.ok) await fail(res, 'Failed to delete entry');
}

export async function getGamePhotos(): Promise<GamePhotosResponse> {
  const res = await fetch(`${API_BASE}/game/photos`);
  if (!res.ok) await fail(res, 'Failed to fetch photos');
  return res.json();
}

export async function getGameRankings(): Promise<GameRankingsResponse> {
  const res = await fetch(`${API_BASE}/game/rankings`);
  if (!res.ok) await fail(res, 'Failed to fetch rankings');
  return res.json();
}

export async function submitGameScore(
  nickname: string,
  timeMs: number,
  gameToken: string,
): Promise<GameScore> {
  const res = await fetch(`${API_BASE}/game/rankings`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ nickname, time_ms: timeMs, game_token: gameToken }),
  });
  if (!res.ok) await fail(res, 'Failed to submit score');
  return res.json();
}

export async function recordGameBeat(): Promise<void> {
  try {
    await fetch(`${API_BASE}/game/beats`, { method: 'POST' });
  } catch {
    // analytics only — never surface
  }
}

export async function getPhotoStorageStatus(): Promise<boolean> {
  const res = await fetch(`${API_BASE}/photos/status`);
  if (!res.ok) return false;
  const data: { available: boolean } = await res.json();
  return data.available;
}

export interface PaginatedPhotos {
  items: PhotoUpload[];
  next_offset: number | null;
}

export async function getPhotos(
  sort: 'recent' | 'popular' = 'recent',
  offset?: number,
  limit = 10,
): Promise<PaginatedPhotos> {
  const params = new URLSearchParams({ sort, limit: String(limit) });
  if (offset !== undefined) params.set('offset', String(offset));
  const res = await fetch(`${API_BASE}/photos?${params}`);
  if (!res.ok) await fail(res, 'Failed to fetch photos');
  return res.json();
}

export async function uploadPhoto(
  name: string,
  image: File,
  password: string,
  original?: File,
): Promise<{ status: string; id: number }> {
  const form = new FormData();
  form.append('name', name);
  form.append('image', image);
  form.append('password', password);
  if (original) {
    form.append('original', original);
  }

  const res = await fetch(`${API_BASE}/photos/upload`, {
    method: 'POST',
    body: form,
  });
  if (!res.ok) await fail(res, 'Failed to upload photo');
  return res.json();
}

export async function verifyPhotoPassword(
  id: number,
  password: string,
): Promise<void> {
  const res = await fetch(`${API_BASE}/photos/${id}/verify`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ password }),
  });
  if (!res.ok) await fail(res, 'Failed to verify password');
}

export async function deleteUserPhoto(
  id: number,
  password: string,
): Promise<void> {
  const res = await fetch(`${API_BASE}/photos/${id}`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ password }),
  });
  if (!res.ok) await fail(res, 'Failed to delete photo');
}

export async function getHallOfFame(): Promise<HallOfFameEntry[]> {
  const res = await fetch(`${API_BASE}/hall-of-fame`);
  if (!res.ok) await fail(res, 'Failed to fetch hall of fame');
  return res.json();
}

export async function createHallOfFameEntry(
  nickname: string,
): Promise<HallOfFameEntry> {
  const res = await fetch(`${API_BASE}/hall-of-fame`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ nickname }),
  });
  if (!res.ok) await fail(res, 'Failed to submit');
  return res.json();
}
