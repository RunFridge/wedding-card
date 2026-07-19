import { ref } from 'vue';
import api from '@/lib/axios';

export interface GameBeat {
  date: string;
  count: number;
}

export function useGameBeats() {
  const beats = ref<GameBeat[]>([]);
  const loading = ref(false);

  async function load() {
    loading.value = true;
    try {
      const { data } = await api.get('/admin/game-beats');
      beats.value = data || [];
    } catch {
      beats.value = [];
    } finally {
      loading.value = false;
    }
  }

  return { beats, loading, load };
}
