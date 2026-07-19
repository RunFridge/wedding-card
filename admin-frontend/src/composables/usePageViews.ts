import { ref } from 'vue';
import api from '@/lib/axios';

export interface PageView {
  date: string;
  count: number;
}

export function usePageViews() {
  const views = ref<PageView[]>([]);
  const loading = ref(false);

  async function load() {
    loading.value = true;
    try {
      const { data } = await api.get('/admin/page-views');
      views.value = data || [];
    } catch {
      views.value = [];
    } finally {
      loading.value = false;
    }
  }

  return { views, loading, load };
}
