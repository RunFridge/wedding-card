import { ref } from 'vue';
import api from '@/lib/axios';
import type { WeddingConfig } from '@/types/admin';

export function useConfig() {
  const config = ref<WeddingConfig | null>(null);
  const loading = ref(false);
  const saving = ref(false);
  const error = ref('');

  async function load() {
    loading.value = true;
    error.value = '';
    try {
      const { data } = await api.get<WeddingConfig>('/admin/config');
      config.value = data;
    } catch {
      error.value = 'Failed to load config';
    } finally {
      loading.value = false;
    }
  }

  async function save(updated: WeddingConfig) {
    saving.value = true;
    error.value = '';
    try {
      await api.put('/admin/config', updated);
      config.value = updated;
    } catch {
      error.value = 'Failed to save config';
      throw error.value;
    } finally {
      saving.value = false;
    }
  }

  return { config, loading, saving, error, load, save };
}
