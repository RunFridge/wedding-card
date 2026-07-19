import { ref } from 'vue';
import api from '@/lib/axios';
import { minLoadingDelay } from '@/lib/utils';
import type { HallOfFameEntry } from '@/types/admin';

export function useHallOfFame() {
  const entries = ref<HallOfFameEntry[]>([]);
  const loading = ref(false);

  async function load() {
    loading.value = true;
    try {
      const [{ data }] = await Promise.all([
        api.get<HallOfFameEntry[]>('/admin/hall-of-fame'),
        minLoadingDelay(),
      ]);
      entries.value = data || [];
    } finally {
      loading.value = false;
    }
  }

  async function deleteEntryRaw(id: number) {
    await api.delete(`/admin/hall-of-fame/${id}`);
  }

  async function deleteEntry(id: number) {
    await deleteEntryRaw(id);
    await load();
  }

  return { entries, loading, load, deleteEntry, deleteEntryRaw };
}
