import { ref } from 'vue';
import api from '@/lib/axios';
import { minLoadingDelay } from '@/lib/utils';
import { i18n } from '@/i18n';
import type { AssetPhoto } from '@/types/admin';

const FULL_MAX_DIM = 2048;
const FULL_QUALITY = 0.85;
const THUMB_MAX_DIM = 400;
const THUMB_QUALITY = 0.85;

function loadImage(file: File): Promise<HTMLImageElement> {
  return new Promise((resolve, reject) => {
    const img = new Image();
    img.onload = () => resolve(img);
    img.onerror = reject;
    img.src = URL.createObjectURL(file);
  });
}

function resizeToBlob(
  img: HTMLImageElement,
  maxDim: number,
  quality: number,
): Promise<Blob> {
  let { naturalWidth: w, naturalHeight: h } = img;
  if (w > maxDim || h > maxDim) {
    if (w > h) {
      h = Math.round((h * maxDim) / w);
      w = maxDim;
    } else {
      w = Math.round((w * maxDim) / h);
      h = maxDim;
    }
  }
  const canvas = document.createElement('canvas');
  canvas.width = w;
  canvas.height = h;
  const ctx = canvas.getContext('2d')!;
  ctx.drawImage(img, 0, 0, w, h);
  return new Promise((resolve, reject) => {
    canvas.toBlob(
      (blob) =>
        blob ? resolve(blob) : reject(new Error('Canvas toBlob failed')),
      'image/jpeg',
      quality,
    );
  });
}

export interface BatchProgress {
  total: number;
  completed: number;
  currentFileName: string;
  errors: { fileName: string; message: string }[];
}

export function useAssetPhotos() {
  const photos = ref<AssetPhoto[]>([]);
  const loading = ref(false);

  async function load() {
    loading.value = true;
    try {
      const [{ data }] = await Promise.all([
        api.get<AssetPhoto[]>('/admin/asset-photos'),
        minLoadingDelay(),
      ]);
      photos.value = data || [];
    } finally {
      loading.value = false;
    }
  }

  async function uploadSingle(file: File, label: string) {
    const img = await loadImage(file);
    const optimized = await resizeToBlob(img, FULL_MAX_DIM, FULL_QUALITY);
    const thumbnail = await resizeToBlob(img, THUMB_MAX_DIM, THUMB_QUALITY);
    URL.revokeObjectURL(img.src);

    const form = new FormData();
    form.append('image', optimized, file.name);
    form.append('thumbnail', thumbnail, `thumb_${file.name}`);
    if (label) form.append('label', label);
    await api.post('/admin/asset-photos', form, {
      headers: { 'Content-Type': undefined },
    });
  }

  async function upload(file: File, label: string) {
    await uploadSingle(file, label);
    await load();
  }

  async function uploadBatch(
    files: File[],
    label: string,
    onProgress?: (progress: BatchProgress) => void,
  ) {
    const progress: BatchProgress = {
      total: files.length,
      completed: 0,
      currentFileName: '',
      errors: [],
    };

    for (const file of files) {
      progress.currentFileName = file.name;
      onProgress?.({ ...progress });
      try {
        await uploadSingle(file, label);
      } catch (err: unknown) {
        const message =
          err instanceof Error
            ? err.message
            : i18n.global.t('assetPhotos.uploadFailed');
        progress.errors.push({ fileName: file.name, message });
      }
      progress.completed++;
    }

    progress.currentFileName = '';
    onProgress?.({ ...progress });
    await load();
    return progress.errors;
  }

  async function toggleGameRaw(id: number, useForGame: boolean) {
    await api.patch(`/admin/asset-photos/${id}/game`, {
      use_for_game: useForGame,
    });
  }

  async function toggleGame(id: number, useForGame: boolean) {
    await toggleGameRaw(id, useForGame);
    await load();
  }

  async function setMain(id: number, isMainPhoto: boolean) {
    await api.patch(`/admin/asset-photos/${id}/main`, {
      is_main_photo: isMainPhoto,
    });
    await load();
  }

  async function updatePhoto(id: number, label: string, sortOrder: number) {
    await api.patch(`/admin/asset-photos/${id}`, {
      label,
      sort_order: sortOrder,
    });
    await load();
  }

  async function deletePhotoRaw(id: number) {
    await api.delete(`/admin/asset-photos/${id}`);
  }

  async function deletePhoto(id: number) {
    await deletePhotoRaw(id);
    await load();
  }

  return {
    photos,
    loading,
    load,
    upload,
    uploadBatch,
    toggleGame,
    toggleGameRaw,
    setMain,
    updatePhoto,
    deletePhoto,
    deletePhotoRaw,
  };
}
