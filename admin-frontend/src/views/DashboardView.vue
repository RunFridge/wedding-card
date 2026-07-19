<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import {
  BookOpen,
  Trophy,
  Camera,
  Shield,
  ListTodo,
  Heart,
  Users,
  Eye,
} from 'lucide-vue-next';
import api from '@/lib/axios';
import { useModeration } from '@/composables/useModeration';
import { useAdminWS } from '@/composables/useAdminWS';
import { usePageViews } from '@/composables/usePageViews';
import { useGameBeats } from '@/composables/useGameBeats';
import VisitChart from '@/components/VisitChart.vue';
import type { ModerationCategoryStats } from '@/types/admin';

const { t } = useI18n();
const { visitorCount, totalHearts } = useAdminWS();
const {
  views: pageViews,
  loading: pvLoading,
  load: loadPageViews,
} = usePageViews();
const {
  beats: gameBeats,
  loading: gbLoading,
  load: loadGameBeats,
} = useGameBeats();

const counts = ref({ guestbook: 0, rankings: 0, photos: 0 });
const totalVisits = computed(() =>
  pageViews.value.reduce((s, v) => s + v.count, 0),
);
const todayVisits = computed(() => {
  const d = new Date();
  const today = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;
  return pageViews.value.find((v) => v.date === today)?.count ?? 0;
});
const avgVisits = computed(() => {
  if (!pageViews.value.length) return 0;
  return Math.round(totalVisits.value / pageViews.value.length);
});

const totalBeats = computed(() =>
  gameBeats.value.reduce((s, v) => s + v.count, 0),
);
const todayBeats = computed(() => {
  const d = new Date();
  const today = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;
  return gameBeats.value.find((v) => v.date === today)?.count ?? 0;
});
const avgBeats = computed(() => {
  if (!gameBeats.value.length) return 0;
  return Math.round(totalBeats.value / gameBeats.value.length);
});
const {
  status: modStatus,
  loading: modLoading,
  error: modError,
  load: loadModeration,
} = useModeration();

let pollTimer: ReturnType<typeof setInterval> | null = null;

onMounted(async () => {
  try {
    const [gb, rk, ph] = await Promise.all([
      api.get('/admin/guestbook'),
      api.get('/admin/game/rankings'),
      api.get('/admin/photos'),
    ]);
    counts.value = {
      guestbook: (gb.data || []).length,
      rankings: (rk.data || []).length,
      photos: (ph.data || []).length,
    };
  } catch {
    // counts stay at 0
  }
  await Promise.all([loadModeration(), loadPageViews(), loadGameBeats()]);
  if (modStatus.value) {
    totalHearts.value = modStatus.value.total_hearts;
  }
  pollTimer = setInterval(loadModeration, 30000);
});

onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer);
});

const cards = [
  {
    key: 'guestbook',
    labelKey: 'dashboard.guestbookEntries',
    icon: BookOpen,
    route: '/guestbook',
  },
  {
    key: 'rankings',
    labelKey: 'dashboard.gameScores',
    icon: Trophy,
    route: '/rankings',
  },
  {
    key: 'photos',
    labelKey: 'dashboard.photoUploads',
    icon: Camera,
    route: '/photos',
  },
] as const;

function statItems(stats: ModerationCategoryStats) {
  return [
    {
      label: t('dashboard.pending'),
      value: stats.pending,
      variant: 'outline' as const,
    },
    {
      label: t('dashboard.approved'),
      value: stats.approved,
      variant: 'default' as const,
    },
    {
      label: t('dashboard.flagged'),
      value: stats.flagged,
      variant: 'destructive' as const,
    },
  ];
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold">{{ t('nav.dashboard') }}</h1>
    <p class="mb-6 text-sm text-muted-foreground">
      {{ t('dashboard.subtitle') }}
    </p>
    <div class="grid gap-4 sm:grid-cols-3">
      <RouterLink
        v-for="card in cards"
        :key="card.key"
        :to="card.route"
        class="block"
      >
        <Card class="transition-shadow hover:shadow-md">
          <CardHeader class="flex flex-row items-center justify-between pb-2">
            <CardTitle class="text-sm font-medium text-muted-foreground">
              {{ t(card.labelKey) }}
            </CardTitle>
            <component :is="card.icon" class="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div class="text-3xl font-bold">
              {{ counts[card.key] }}
            </div>
          </CardContent>
        </Card>
      </RouterLink>
    </div>

    <!-- Page Views -->
    <h2 class="mb-4 mt-8 text-xl font-semibold">
      {{ t('dashboard.pageViews') }}
    </h2>
    <div class="grid gap-4 sm:grid-cols-3 mb-4">
      <Card>
        <CardHeader class="flex flex-row items-center justify-between pb-2">
          <CardTitle class="text-sm font-medium text-muted-foreground">{{
            t('dashboard.today')
          }}</CardTitle>
          <Eye class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-3xl font-bold">
            {{ todayVisits.toLocaleString() }}
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader class="flex flex-row items-center justify-between pb-2">
          <CardTitle class="text-sm font-medium text-muted-foreground">{{
            t('dashboard.dailyAverage')
          }}</CardTitle>
          <Eye class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-3xl font-bold">{{ avgVisits.toLocaleString() }}</div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader class="flex flex-row items-center justify-between pb-2">
          <CardTitle class="text-sm font-medium text-muted-foreground">{{
            t('dashboard.total')
          }}</CardTitle>
          <Eye class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-3xl font-bold">
            {{ totalVisits.toLocaleString() }}
          </div>
        </CardContent>
      </Card>
    </div>
    <Card>
      <CardHeader class="pb-2">
        <CardTitle class="text-sm font-medium text-muted-foreground">
          {{ t('dashboard.dailyVisits') }}
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div
          v-if="pvLoading && !pageViews.length"
          class="text-sm text-muted-foreground py-8 text-center"
        >
          {{ t('common.loading') }}
        </div>
        <div
          v-else-if="!pageViews.length"
          class="text-sm text-muted-foreground py-8 text-center"
        >
          {{ t('dashboard.noVisitData') }}
        </div>
        <VisitChart v-else :data="pageViews" />
      </CardContent>
    </Card>

    <!-- Game Beats -->
    <h2 class="mb-4 mt-8 text-xl font-semibold">
      {{ t('dashboard.gameBeats') }}
    </h2>
    <div class="grid gap-4 sm:grid-cols-3 mb-4">
      <Card>
        <CardHeader class="flex flex-row items-center justify-between pb-2">
          <CardTitle class="text-sm font-medium text-muted-foreground">{{
            t('dashboard.today')
          }}</CardTitle>
          <Trophy class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-3xl font-bold">
            {{ todayBeats.toLocaleString() }}
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader class="flex flex-row items-center justify-between pb-2">
          <CardTitle class="text-sm font-medium text-muted-foreground">{{
            t('dashboard.dailyAverage')
          }}</CardTitle>
          <Trophy class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-3xl font-bold">{{ avgBeats.toLocaleString() }}</div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader class="flex flex-row items-center justify-between pb-2">
          <CardTitle class="text-sm font-medium text-muted-foreground">{{
            t('dashboard.total')
          }}</CardTitle>
          <Trophy class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-3xl font-bold">
            {{ totalBeats.toLocaleString() }}
          </div>
        </CardContent>
      </Card>
    </div>
    <Card>
      <CardHeader class="pb-2">
        <CardTitle class="text-sm font-medium text-muted-foreground">
          {{ t('dashboard.dailyGameBeats') }}
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div
          v-if="gbLoading && !gameBeats.length"
          class="text-sm text-muted-foreground py-8 text-center"
        >
          {{ t('common.loading') }}
        </div>
        <div
          v-else-if="!gameBeats.length"
          class="text-sm text-muted-foreground py-8 text-center"
        >
          {{ t('dashboard.noGameBeatData') }}
        </div>
        <VisitChart v-else :data="gameBeats" />
      </CardContent>
    </Card>

    <!-- Hearts (always visible) -->
    <h2 class="mb-4 mt-8 text-xl font-semibold">{{ t('dashboard.hearts') }}</h2>

    <div v-if="modLoading && !modStatus" class="text-sm text-muted-foreground">
      {{ t('common.loading') }}
    </div>

    <div
      v-else-if="modError && !modStatus"
      class="text-sm text-muted-foreground"
    >
      {{ t('dashboard.unableStatus') }}
    </div>

    <div v-else-if="modStatus" class="grid gap-4 sm:grid-cols-2">
      <Card>
        <CardHeader class="flex flex-row items-center justify-between pb-2">
          <CardTitle class="text-sm font-medium text-muted-foreground">
            {{ t('dashboard.totalHearts') }}
          </CardTitle>
          <Heart class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-3xl font-bold">{{ totalHearts }}</div>
          <p class="mt-1 text-xs text-muted-foreground">
            {{ t('dashboard.acrossAllPhotos') }}
          </p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between pb-2">
          <CardTitle class="text-sm font-medium text-muted-foreground">
            {{ t('dashboard.onlineVisitors') }}
          </CardTitle>
          <Users class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-3xl font-bold">{{ visitorCount }}</div>
          <p class="mt-1 text-xs text-muted-foreground">
            {{ t('dashboard.liveConnections') }}
          </p>
        </CardContent>
      </Card>
    </div>

    <!-- Moderation Status -->
    <h2 class="mb-4 mt-8 text-xl font-semibold">
      {{ t('dashboard.moderationStatus') }}
    </h2>

    <div v-if="modLoading && !modStatus" class="text-sm text-muted-foreground">
      {{ t('dashboard.loadingModerationStatus') }}
    </div>

    <div
      v-else-if="modError && !modStatus"
      class="text-sm text-muted-foreground"
    >
      {{ t('dashboard.unableModerationStatus') }}
    </div>

    <div
      v-else-if="modStatus && !modStatus.enabled"
      class="text-sm text-muted-foreground"
    >
      {{ t('dashboard.notEnabled') }}
    </div>

    <div v-else-if="modStatus" class="grid gap-4 sm:grid-cols-3">
      <Card>
        <CardHeader class="flex flex-row items-center justify-between pb-2">
          <CardTitle class="text-sm font-medium text-muted-foreground">
            {{ t('dashboard.queueDepth') }}
          </CardTitle>
          <ListTodo class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-3xl font-bold">{{ modStatus.queue_depth }}</div>
          <p class="mt-1 text-xs text-muted-foreground">
            {{ t('dashboard.jobsWaiting') }}
          </p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between pb-2">
          <CardTitle class="text-sm font-medium text-muted-foreground">
            {{ t('nav.guestbook') }}
          </CardTitle>
          <Shield class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-3xl font-bold">{{ modStatus.guestbook.total }}</div>
          <div class="mt-2 flex flex-wrap gap-1.5">
            <Badge
              v-for="item in statItems(modStatus.guestbook)"
              :key="item.label"
              :variant="item.variant"
            >
              {{ item.label }}: {{ item.value }}
            </Badge>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between pb-2">
          <CardTitle class="text-sm font-medium text-muted-foreground">
            {{ t('nav.photos') }}
          </CardTitle>
          <Shield class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-3xl font-bold">{{ modStatus.photos.total }}</div>
          <div class="mt-2 flex flex-wrap gap-1.5">
            <Badge
              v-for="item in statItems(modStatus.photos)"
              :key="item.label"
              :variant="item.variant"
            >
              {{ item.label }}: {{ item.value }}
            </Badge>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
