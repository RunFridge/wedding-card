import { ref } from 'vue';
import api from '@/lib/axios';
import { minLoadingDelay } from '@/lib/utils';
import type { PhotoUpload } from '@/types/admin';

export function usePhotos() {
  const photos = ref<PhotoUpload[]>([]);
  const loading = ref(false);

  async function load() {
    loading.value = true;
    try {
      const [{ data }] = await Promise.all([
        api.get<PhotoUpload[]>('/admin/photos'),
        minLoadingDelay(),
      ]);
      photos.value = data || [];
    } finally {
      loading.value = false;
    }
  }

  async function toggleVisibilityRaw(id: number, hidden: boolean) {
    await api.patch(`/admin/photos/${id}/visibility`, { hidden });
  }

  async function toggleVisibility(id: number, hidden: boolean) {
    await toggleVisibilityRaw(id, hidden);
    await load();
  }

  async function deletePhotoRaw(id: number) {
    await api.delete(`/admin/photos/${id}`);
  }

  async function deletePhoto(id: number) {
    await deletePhotoRaw(id);
    await load();
  }

  async function resetHearts(id: number) {
    await api.post(`/admin/photos/${id}/reset-hearts`);
    await load();
  }

  return {
    photos,
    loading,
    load,
    toggleVisibility,
    toggleVisibilityRaw,
    deletePhoto,
    deletePhotoRaw,
    resetHearts,
  };
}
