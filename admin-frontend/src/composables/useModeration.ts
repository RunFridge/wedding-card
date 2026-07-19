import { ref } from 'vue';
import api from '@/lib/axios';
import type { ModerationStatus } from '@/types/admin';

export function useModeration() {
  const status = ref<ModerationStatus | null>(null);
  const loading = ref(false);
  const error = ref(false);

  async function load() {
    loading.value = true;
    error.value = false;
    try {
      const { data } = await api.get<ModerationStatus>(
        '/admin/moderation/status',
      );
      status.value = data;
    } catch {
      error.value = true;
    } finally {
      loading.value = false;
    }
  }

  return { status, loading, error, load };
}
