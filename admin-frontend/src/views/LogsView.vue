<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Trash2, ArrowDownToLine } from 'lucide-vue-next';

const { t } = useI18n();

const lines = ref<string[]>([]);
const connected = ref(false);
const autoScroll = ref(true);
const container = ref<HTMLElement | null>(null);
let eventSource: EventSource | null = null;

function connect() {
  eventSource = new EventSource('/api/admin/logs');

  eventSource.onopen = () => {
    connected.value = true;
  };

  eventSource.onmessage = (e) => {
    lines.value.push(e.data);
    if (lines.value.length > 2000) {
      lines.value = lines.value.slice(-1500);
    }
  };

  eventSource.onerror = () => {
    connected.value = false;
    eventSource?.close();
    setTimeout(connect, 3000);
  };
}

interface LogPart {
  text: string;
  cls: string;
}

function parseLine(line: string): LogPart[] {
  // Format: "2026/03/25 13:30:45 192.168.1.5 GET /api/config 200 2.1ms 1234 bytes"
  const match = line.match(
    /^(\d{4}\/\d{2}\/\d{2}\s\d{2}:\d{2}:\d{2})\s+(\S+)\s+(GET|POST|PUT|PATCH|DELETE|HEAD|OPTIONS)\s+(\S+)\s+(\d{3})\s+(\S+)\s+(.*)$/,
  );
  if (!match) return [{ text: line, cls: 'text-zinc-400' }];

  const [, timestamp, ip, method, path, status, duration, rest] = match;
  const statusNum = parseInt(status);

  const methodCls =
    method === 'GET'
      ? 'text-green-400'
      : method === 'POST'
        ? 'text-yellow-400'
        : method === 'DELETE'
          ? 'text-red-400'
          : 'text-orange-400';

  const statusCls =
    statusNum < 300
      ? 'text-green-400'
      : statusNum < 400
        ? 'text-blue-400'
        : statusNum < 500
          ? 'text-yellow-400'
          : 'text-red-400';

  return [
    { text: timestamp, cls: 'text-zinc-500' },
    { text: ' ' + ip, cls: 'text-cyan-400' },
    { text: ' ' + method, cls: methodCls + ' font-bold' },
    { text: ' ' + path, cls: 'text-zinc-200' },
    { text: ' ' + status, cls: statusCls + ' font-bold' },
    { text: ' ' + duration, cls: 'text-zinc-500' },
    { text: ' ' + rest, cls: 'text-zinc-600' },
  ];
}

function clearLogs() {
  lines.value = [];
}

function scrollToBottom() {
  if (container.value) {
    container.value.scrollTop = container.value.scrollHeight;
  }
}

watch(
  () => lines.value.length,
  () => {
    if (autoScroll.value) {
      nextTick(scrollToBottom);
    }
  },
);

function onScroll() {
  if (!container.value) return;
  const { scrollTop, scrollHeight, clientHeight } = container.value;
  autoScroll.value = scrollHeight - scrollTop - clientHeight < 40;
}

onMounted(connect);

onUnmounted(() => {
  eventSource?.close();
  eventSource = null;
});
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-4">
      <div>
        <h1 class="text-2xl font-bold">{{ t('logs.title') }}</h1>
        <p class="text-sm text-muted-foreground">{{ t('logs.subtitle') }}</p>
      </div>
      <div class="flex items-center gap-2">
        <span
          class="inline-flex items-center gap-1.5 text-xs"
          :class="connected ? 'text-green-600' : 'text-muted-foreground'"
        >
          <span
            class="w-2 h-2 rounded-full"
            :class="connected ? 'bg-green-500' : 'bg-muted-foreground'"
          />
          {{ connected ? t('logs.connected') : t('logs.reconnecting') }}
        </span>
      </div>
    </div>

    <Card>
      <CardHeader class="flex flex-row items-center justify-between pb-2">
        <CardTitle class="text-sm font-medium text-muted-foreground">
          {{ t('logs.lineCount', lines.length) }}
        </CardTitle>
        <div class="flex items-center gap-1">
          <button
            class="p-1.5 rounded hover:bg-muted transition-colors"
            :title="t('logs.scrollToBottom')"
            @click="
              autoScroll = true;
              scrollToBottom();
            "
          >
            <ArrowDownToLine class="w-4 h-4 text-muted-foreground" />
          </button>
          <button
            class="p-1.5 rounded hover:bg-muted transition-colors"
            :title="t('logs.clear')"
            @click="clearLogs"
          >
            <Trash2 class="w-4 h-4 text-muted-foreground" />
          </button>
        </div>
      </CardHeader>
      <CardContent>
        <div
          ref="container"
          class="h-[60vh] overflow-y-auto bg-zinc-950 rounded p-3 font-mono text-xs leading-relaxed"
          @scroll="onScroll"
        >
          <div v-if="!lines.length" class="text-zinc-500 text-center py-8">
            {{ t('logs.waiting') }}
          </div>
          <div
            v-for="(line, i) in lines"
            :key="i"
            class="whitespace-pre-wrap break-all"
          >
            <span
              v-for="(part, j) in parseLine(line)"
              :key="j"
              :class="part.cls"
              >{{ part.text }}</span
            >
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
