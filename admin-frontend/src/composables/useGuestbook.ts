import { ref } from 'vue';
import api from '@/lib/axios';
import { minLoadingDelay } from '@/lib/utils';
import type { GuestbookEntry } from '@/types/admin';

export function useGuestbook() {
  const entries = ref<GuestbookEntry[]>([]);
  const loading = ref(false);

  async function load() {
    loading.value = true;
    try {
      const [{ data }] = await Promise.all([
        api.get<GuestbookEntry[]>('/admin/guestbook'),
        minLoadingDelay(),
      ]);
      entries.value = data || [];
    } finally {
      loading.value = false;
    }
  }

  async function toggleVisibilityRaw(id: number, hidden: boolean) {
    await api.patch(`/admin/guestbook/${id}/visibility`, { hidden });
  }

  async function toggleVisibility(id: number, hidden: boolean) {
    await toggleVisibilityRaw(id, hidden);
    await load();
  }

  async function deleteEntryRaw(id: number) {
    await api.delete(`/admin/guestbook/${id}`);
  }

  async function deleteEntry(id: number) {
    await deleteEntryRaw(id);
    await load();
  }

  return {
    entries,
    loading,
    load,
    toggleVisibility,
    toggleVisibilityRaw,
    deleteEntry,
    deleteEntryRaw,
  };
}
