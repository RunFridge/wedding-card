import { ref } from 'vue';
import api from '@/lib/axios';
import { minLoadingDelay } from '@/lib/utils';
import type { GameScore } from '@/types/admin';

export function useRankings() {
  const scores = ref<GameScore[]>([]);
  const loading = ref(false);

  async function load() {
    loading.value = true;
    try {
      const [{ data }] = await Promise.all([
        api.get<GameScore[]>('/admin/game/rankings'),
        minLoadingDelay(),
      ]);
      scores.value = data || [];
    } finally {
      loading.value = false;
    }
  }

  async function deleteScoreRaw(id: number) {
    await api.delete(`/admin/game/rankings/${id}`);
  }

  async function deleteScore(id: number) {
    await deleteScoreRaw(id);
    await load();
  }

  async function purgeAll() {
    await api.post('/admin/game/rankings/purge');
    await load();
  }

  return { scores, loading, load, deleteScore, deleteScoreRaw, purgeAll };
}
