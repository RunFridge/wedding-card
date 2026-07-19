import { reactive, ref } from 'vue';
import type { PhotoUpload } from '../types';
import { clearAll as clearAllStore } from '../lib/store';

const API_BASE = '/api';

const heartCounts = reactive<Record<number, number>>({});
const connected = ref(false);

let ws: WebSocket | null = null;
let reconnectTimer: ReturnType<typeof setTimeout> | null = null;
let destroyed = false;

const pendingHearts: Record<number, number> = {};
let flushTimer: ReturnType<typeof setTimeout> | null = null;

const contentCallbacks: Record<string, Set<() => void>> = {};
const rankingCallbacks = new Set<(nickname: string, timeMs: number) => void>();
const gameResetCallbacks = new Set<() => void>();

function initHearts(photos: PhotoUpload[]) {
  for (const p of photos) {
    heartCounts[p.id] = p.hearts;
  }
}

function connect() {
  if (ws) return;
  destroyed = false;
  const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:';
  ws = new WebSocket(`${protocol}//${location.host}${API_BASE}/ws`);

  ws.onopen = () => {
    connected.value = true;
  };

  ws.onclose = () => {
    connected.value = false;
    ws = null;
    if (!destroyed) {
      reconnectTimer = setTimeout(connect, 3000);
    }
  };

  ws.onerror = () => {
    ws?.close();
  };

  ws.onmessage = (event) => {
    try {
      const msg = JSON.parse(event.data);
      if (msg.type === 'hearts_update' && Array.isArray(msg.updates)) {
        for (const u of msg.updates) {
          heartCounts[u.photo_id] = u.hearts;
        }
      } else if (msg.type === 'heart_increment' && msg.photo_id) {
        heartCounts[msg.photo_id] =
          (heartCounts[msg.photo_id] || 0) + msg.count;
      } else if (msg.type === 'content_update' && msg.content_type) {
        const cbs = contentCallbacks[msg.content_type];
        if (cbs) {
          for (const cb of cbs) cb();
        }
      } else if (msg.type === 'ranking_update' && msg.nickname) {
        for (const cb of rankingCallbacks) cb(msg.nickname, msg.time_ms);
      } else if (msg.type === 'game_reset') {
        clearAllStore();
        for (const cb of gameResetCallbacks) cb();
      }
    } catch {
      // ignore malformed messages
    }
  };
}

function flushPending() {
  flushTimer = null;
  const entries = Object.entries(pendingHearts);
  if (entries.length === 0) return;

  for (const [idStr, count] of entries) {
    const id = Number(idStr);
    delete pendingHearts[id];
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ type: 'heart', photo_id: id, count }));
    }
  }
}

function sendHeart(photoId: number) {
  heartCounts[photoId] = (heartCounts[photoId] || 0) + 1;
  navigator.vibrate?.(10);
  pendingHearts[photoId] = (pendingHearts[photoId] || 0) + 1;

  if (!flushTimer) {
    flushTimer = setTimeout(flushPending, 300);
  }
}

function disconnect() {
  destroyed = true;
  if (flushTimer) {
    clearTimeout(flushTimer);
    flushPending();
  }
  if (reconnectTimer) clearTimeout(reconnectTimer);
  ws?.close();
  ws = null;
}

function onContentUpdate(contentType: string, callback: () => void) {
  if (!contentCallbacks[contentType]) {
    contentCallbacks[contentType] = new Set();
  }
  contentCallbacks[contentType].add(callback);
}

function offContentUpdate(contentType: string, callback: () => void) {
  contentCallbacks[contentType]?.delete(callback);
}

function onRankingUpdate(callback: (nickname: string, timeMs: number) => void) {
  rankingCallbacks.add(callback);
}

function offRankingUpdate(
  callback: (nickname: string, timeMs: number) => void,
) {
  rankingCallbacks.delete(callback);
}

function onGameReset(callback: () => void) {
  gameResetCallbacks.add(callback);
}

function offGameReset(callback: () => void) {
  gameResetCallbacks.delete(callback);
}

export function useHearts() {
  return {
    heartCounts,
    connected,
    connect,
    initHearts,
    sendHeart,
    disconnect,
    onContentUpdate,
    offContentUpdate,
    onRankingUpdate,
    offRankingUpdate,
    onGameReset,
    offGameReset,
  };
}
