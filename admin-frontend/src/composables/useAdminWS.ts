import { ref } from 'vue';

const visitorCount = ref(0);
const totalHearts = ref(0);
const connected = ref(false);

let ws: WebSocket | null = null;
let reconnectTimer: ReturnType<typeof setTimeout> | null = null;
let destroyed = false;
let consecutiveFailures = 0;
const MAX_RECONNECT_FAILURES = 5;

function connect() {
  if (ws) return;
  destroyed = false;
  const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:';
  ws = new WebSocket(`${protocol}//${location.host}/api/admin/ws`);

  ws.onopen = () => {
    connected.value = true;
    consecutiveFailures = 0;
  };

  ws.onclose = (event) => {
    connected.value = false;
    ws = null;
    if (destroyed) return;

    // 1006 = abnormal closure (typical for HTTP 401 on upgrade)
    if (event.code === 1006) {
      consecutiveFailures++;
    }

    if (consecutiveFailures >= MAX_RECONNECT_FAILURES) {
      window.location.href = '/-/admin/login';
      return;
    }

    const delay = Math.min(3000 * 2 ** (consecutiveFailures - 1), 30000);
    reconnectTimer = setTimeout(connect, delay);
  };

  ws.onerror = () => {
    ws?.close();
  };

  ws.onmessage = (event) => {
    try {
      const msg = JSON.parse(event.data);
      if (msg.type === 'presence') {
        visitorCount.value = msg.visitor_count;
      } else if (msg.type === 'heart_increment') {
        totalHearts.value += msg.count;
      }
    } catch {
      // ignore malformed messages
    }
  };
}

function disconnect() {
  destroyed = true;
  consecutiveFailures = 0;
  if (reconnectTimer) clearTimeout(reconnectTimer);
  ws?.close();
  ws = null;
}

export function useAdminWS() {
  return { visitorCount, totalHearts, connected, connect, disconnect };
}
