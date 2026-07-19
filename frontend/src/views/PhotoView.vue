<template>
  <div class="max-w-md mx-auto px-4 py-6">
    <div class="wooden-panel p-2 mb-6">
      <div class="wooden-panel-inner px-4 py-3 text-center">
        <h1 class="text-xl text-secondary mb-1">{{ t('photo.title') }}</h1>
        <p class="text-wood-dark text-sm break-keep">
          {{ t('photo.subtitle') }}
        </p>
      </div>
    </div>

    <!-- Time gate: before 1 hour prior to wedding -->
    <div
      v-if="!timeGateOpen"
      class="parchment-bg p-6 text-center text-wood-dark/80 mb-6"
    >
      <p class="text-base mb-3 break-keep">
        {{ t('photo.availableFrom', { duration: gateDuration }) }}
      </p>
      <div class="flex justify-center gap-3 text-wood-dark">
        <div class="flex flex-col items-center">
          <span class="text-2xl font-bold text-secondary">{{
            countdown.days
          }}</span>
          <span class="text-[10px] text-wood-dark/50">{{
            t('photo.unitDays')
          }}</span>
        </div>
        <span class="text-2xl text-wood-dark/30">:</span>
        <div class="flex flex-col items-center">
          <span class="text-2xl font-bold text-secondary">{{
            countdown.hours
          }}</span>
          <span class="text-[10px] text-wood-dark/50">{{
            t('photo.unitHours')
          }}</span>
        </div>
        <span class="text-2xl text-wood-dark/30">:</span>
        <div class="flex flex-col items-center">
          <span class="text-2xl font-bold text-secondary">{{
            countdown.minutes
          }}</span>
          <span class="text-[10px] text-wood-dark/50">{{
            t('photo.unitMinutes')
          }}</span>
        </div>
        <span class="text-2xl text-wood-dark/30">:</span>
        <div class="flex flex-col items-center">
          <span class="text-2xl font-bold text-secondary">{{
            countdown.seconds
          }}</span>
          <span class="text-[10px] text-wood-dark/50">{{
            t('photo.unitSeconds')
          }}</span>
        </div>
      </div>
    </div>

    <!-- Upload not available -->
    <div v-else-if="checkingStatus" class="flex justify-center py-6">
      <div class="pixel-spinner"></div>
    </div>
    <div
      v-else-if="!uploadAvailable"
      class="parchment-bg p-6 text-center text-wood-dark/80 mb-6"
    >
      <p class="text-base">{{ t('photo.notReady') }}</p>
      <p class="text-sm mt-2 text-wood-dark/50">{{ t('photo.tryLater') }}</p>
    </div>

    <!-- Mobile: upload form -->
    <div v-else class="wooden-panel p-2 mb-6">
      <div class="wooden-panel-inner p-4">
        <form @submit.prevent="submitPhoto">
          <div class="mb-3">
            <label class="block text-secondary text-sm mb-1">{{
              t('photo.nameLabel')
            }}</label>
            <input
              v-model="name"
              type="text"
              maxlength="50"
              required
              class="pixel-input w-full"
              :placeholder="t('common.namePlaceholder')"
            />
          </div>
          <div class="mb-3">
            <label class="block text-secondary text-sm mb-1">{{
              t('photo.passwordLabel')
            }}</label>
            <input
              v-model="uploadPassword"
              type="password"
              minlength="5"
              maxlength="30"
              required
              class="pixel-input w-full"
              :placeholder="t('photo.passwordPlaceholder')"
            />
          </div>
          <div class="mb-3">
            <input
              ref="fileInput"
              type="file"
              accept="image/*"
              capture="environment"
              class="hidden"
              @change="onFileChange"
            />
            <input
              ref="galleryInput"
              type="file"
              accept="image/*"
              multiple
              class="hidden"
              @change="onFileChange"
            />
            <div v-if="queue.length === 0" class="flex gap-2">
              <button
                type="button"
                :disabled="uploading"
                class="pixel-btn flex-1 flex items-center justify-center gap-2 whitespace-nowrap sm:hidden"
                @click="fileInput?.click()"
              >
                <TwEmoji emoji="📷" size="1.25rem" />
                <span>{{ t('photo.openCamera') }}</span>
              </button>
              <button
                type="button"
                :disabled="uploading"
                class="pixel-btn flex-1 flex items-center justify-center gap-2 whitespace-nowrap"
                @click="galleryInput?.click()"
              >
                <TwEmoji emoji="🖼️" size="1.25rem" />
                <span>{{ t('photo.choosePhoto') }}</span>
              </button>
            </div>
          </div>
          <div v-if="queue.length > 0" class="mb-3">
            <p class="text-xs text-wood-dark/60 mb-2">
              {{ t('photo.editHint') }}
            </p>
            <div class="grid grid-cols-3 gap-2">
              <div v-for="(item, i) in queue" :key="item.preview" class="relative">
                <img
                  :src="item.preview"
                  :alt="t('photo.previewAlt')"
                  class="w-full aspect-square object-cover rounded border-2 cursor-pointer"
                  :class="
                    item.status === 'error'
                      ? 'border-red-500'
                      : 'border-wood-dark/20'
                  "
                  @click="!uploading && openEditor(i)"
                />
                <div
                  v-if="item.status === 'uploading'"
                  class="absolute inset-0 flex items-center justify-center bg-black/30 rounded"
                >
                  <div class="pixel-spinner"></div>
                </div>
                <button
                  v-if="!uploading"
                  type="button"
                  :aria-label="t('photo.removePhotoAria')"
                  class="absolute -top-1.5 -right-1.5 w-6 h-6 rounded-full bg-wood-dark text-white text-sm leading-none flex items-center justify-center"
                  @click="removeItem(i)"
                >
                  ×
                </button>
              </div>
            </div>
          </div>
          <div v-if="queue.length > 0" class="flex gap-2 mb-3">
            <button
              v-if="!uploading"
              type="button"
              :aria-label="t('photo.openCamera')"
              class="pixel-btn flex-none px-3 sm:hidden"
              @click="fileInput?.click()"
            >
              <TwEmoji emoji="📷" size="1.25rem" />
            </button>
            <button
              v-if="!uploading"
              type="button"
              class="pixel-btn flex-1 whitespace-nowrap"
              @click="galleryInput?.click()"
            >
              {{ t('photo.addMore') }}
            </button>
            <button
              type="submit"
              :disabled="uploading"
              class="pixel-btn flex-1 min-w-0 truncate"
            >
              {{
                uploading
                  ? t('photo.uploadingProgress', {
                      done: doneCount,
                      total: queue.length,
                    })
                  : t('photo.uploadCount', { n: queue.length }, queue.length)
              }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Sort toggle -->
    <div v-if="photos.length > 0" class="flex justify-center gap-2 mb-4">
      <button
        class="sort-btn text-xs px-4 py-1"
        :class="sortMode === 'recent' ? 'sort-btn-active' : 'sort-btn-inactive'"
        :aria-pressed="sortMode === 'recent'"
        @click="switchSort('recent')"
      >
        {{ t('photo.sortRecent') }}
      </button>
      <button
        class="sort-btn text-xs px-4 py-1"
        :class="
          sortMode === 'popular' ? 'sort-btn-active' : 'sort-btn-inactive'
        "
        :aria-pressed="sortMode === 'popular'"
        @click="switchSort('popular')"
      >
        {{ t('photo.sortPopular') }}
      </button>
    </div>

    <!-- Photo gallery -->
    <div v-if="loadingPhotos" class="flex justify-center py-6">
      <div class="pixel-spinner"></div>
    </div>

    <div
      v-else-if="photos.length === 0"
      class="parchment-bg p-6 text-center text-wood-dark/80"
    >
      {{ t('photo.empty') }}
    </div>

    <div v-else class="masonry-grid">
      <div v-for="(col, ci) in photoColumns" :key="ci" class="masonry-col">
        <div
          v-for="photo in col"
          :key="photo.id"
          class="parchment-bg p-2 mb-2 relative"
        >
          <div class="absolute top-3 right-3 z-10">
            <button
              :aria-label="t('photo.menuAria')"
              aria-haspopup="menu"
              :aria-expanded="openKebab === photo.id"
              class="w-8 h-8 rounded-full text-wood-dark/80 hover:bg-wood-dark/20 flex items-center justify-center text-lg font-bold leading-none transition-colors"
              @click.stop="toggleKebab(photo.id)"
            >
              ⋮
            </button>
            <div
              v-if="openKebab === photo.id"
              role="menu"
              class="absolute right-0 top-full mt-1 bg-parchment border border-wood-dark/20 shadow-md rounded z-10 whitespace-nowrap"
            >
              <button
                role="menuitem"
                class="block w-full text-left text-xs text-red-600 hover:bg-wood-dark/10 px-3 py-1.5"
                @click="
                  startPhotoDelete(photo);
                  openKebab = null;
                "
              >
                {{ t('photo.deletePhoto') }}
              </button>
            </div>
          </div>
          <img
            :src="photo.thumbnail || photo.url"
            :alt="photo.name"
            class="w-full rounded border border-wood-dark/10 cursor-pointer"
            @click="lightboxIndex = photos.indexOf(photo)"
          />
          <div class="flex justify-between items-center mt-1">
            <p class="text-sm text-wood-dark/80 truncate flex-1">
              {{ photo.name }}
            </p>
            <button
              class="heart-btn text-sm px-2 py-1 flex items-center gap-1"
              :aria-label="t('common.like')"
              @click="onHeartClick(photo, $event)"
            >
              <span class="heart-icon">♥</span>
              <span class="text-wood-dark/50">{{
                heartCounts[photo.id] || 0
              }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Infinite scroll sentinel -->
    <div
      v-if="nextOffset !== null"
      ref="photoSentinelRef"
      class="flex justify-center py-4"
    >
      <div v-if="loadingMore" class="pixel-spinner"></div>
    </div>

    <!-- Phaser heart overlay -->
    <Teleport to="body">
      <div
        ref="heartOverlayRef"
        class="fixed inset-0 pointer-events-none"
        style="z-index: 9999"
      ></div>
    </Teleport>

    <PhotoEditor
      :visible="editorVisible"
      :source-url="editingPreview"
      @update:visible="editorVisible = $event"
      @confirm="onEditorConfirm"
    />

    <LightboxViewer
      v-model="lightboxIndex"
      :items="lightboxItems"
      :heart-counts="heartCounts"
      @heart="onLightboxHeart"
    />

    <StardewDialog
      :visible="alertDialog.visible"
      :message="alertDialog.message"
      @close="alertDialog.visible = false"
    />

    <StardewDialog
      :visible="photoPasswordDialog.visible"
      :message="photoPasswordDialog.message"
      input-mode="password"
      @close="photoPasswordDialog.visible = false"
      @cancel="photoPasswordDialog.visible = false"
      @confirm="onPhotoDeleteConfirm"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, onMounted, onUnmounted } from 'vue';
import { onBeforeRouteLeave } from 'vue-router';
import { useI18n } from 'vue-i18n';
import Phaser from 'phaser';
import {
  getPhotos,
  getPhotoStorageStatus,
  uploadPhoto,
  deleteUserPhoto,
  ApiError,
  apiErrorMessage,
} from '../services/api';
import {
  WEDDING_DATETIME,
  PHOTO_UPLOAD_ENABLED,
  PHOTO_UPLOAD_HOURS_BEFORE,
} from '../config/wedding';
import {
  compressImage,
  runWithConcurrency,
  UPLOAD_CONCURRENCY,
} from '../lib/upload';
import { useHearts } from '../composables/useHearts';
import { useAchievements } from '../composables/useAchievements';
import HeartScene from '../game/HeartScene';
import StardewDialog from '../components/StardewDialog.vue';
import LightboxViewer from '../components/LightboxViewer.vue';
import PhotoEditor from '../components/PhotoEditor.vue';
import TwEmoji from '../components/TwEmoji.vue';
import type { PhotoUpload, DialogState } from '../types';

const { t } = useI18n();

const OPEN_TIME = new Date(
  WEDDING_DATETIME.getTime() - PHOTO_UPLOAD_HOURS_BEFORE * 60 * 60 * 1000,
);

const gateDuration = computed(() => {
  if (PHOTO_UPLOAD_HOURS_BEFORE >= 24) {
    const days = Math.round(PHOTO_UPLOAD_HOURS_BEFORE / 24);
    return t('photo.durationDays', { n: days }, days);
  }
  return t(
    'photo.durationHours',
    { n: PHOTO_UPLOAD_HOURS_BEFORE },
    PHOTO_UPLOAD_HOURS_BEFORE,
  );
});

const timeGateOpen = ref(
  PHOTO_UPLOAD_ENABLED || Date.now() >= OPEN_TIME.getTime(),
);
const countdown = reactive({
  days: '00',
  hours: '00',
  minutes: '00',
  seconds: '00',
});
let countdownTimer: ReturnType<typeof setInterval> | null = null;

function updateCountdown() {
  if (PHOTO_UPLOAD_ENABLED) {
    timeGateOpen.value = true;
    if (countdownTimer) {
      clearInterval(countdownTimer);
      countdownTimer = null;
    }
    checkStorageStatus();
    return;
  }
  const diff = OPEN_TIME.getTime() - Date.now();
  if (diff <= 0) {
    timeGateOpen.value = true;
    if (countdownTimer) {
      clearInterval(countdownTimer);
      countdownTimer = null;
    }
    checkStorageStatus();
    return;
  }
  const d = Math.floor(diff / (1000 * 60 * 60 * 24));
  const h = Math.floor((diff / (1000 * 60 * 60)) % 24);
  const m = Math.floor((diff / (1000 * 60)) % 60);
  const s = Math.floor((diff / 1000) % 60);
  countdown.days = String(d).padStart(2, '0');
  countdown.hours = String(h).padStart(2, '0');
  countdown.minutes = String(m).padStart(2, '0');
  countdown.seconds = String(s).padStart(2, '0');
}

const uploadAvailable = ref(true);
const checkingStatus = ref(true);

interface QueueItem {
  file: File;
  original: File;
  preview: string;
  status: 'pending' | 'uploading' | 'done' | 'error';
}

const name = ref('');
const uploadPassword = ref('');
const queue = ref<QueueItem[]>([]);
const uploading = ref(false);
const fileInput = ref<HTMLInputElement | null>(null);
const galleryInput = ref<HTMLInputElement | null>(null);
const editorVisible = ref(false);
const editingIndex = ref<number | null>(null);

const doneCount = computed(
  () => queue.value.filter((item) => item.status === 'done').length,
);
const editingPreview = computed(() =>
  editingIndex.value !== null
    ? (queue.value[editingIndex.value]?.preview ?? '')
    : '',
);

const photos = ref<PhotoUpload[]>([]);
const photoColumns = computed(() => {
  const left: PhotoUpload[] = [];
  const right: PhotoUpload[] = [];
  photos.value.forEach((p, i) => (i % 2 === 0 ? left : right).push(p));
  return [left, right];
});
const loadingPhotos = ref(true);
const loadingMore = ref(false);
const nextOffset = ref<number | null>(null);
const photoSentinelRef = ref<HTMLElement | null>(null);
let photoScrollObserver: IntersectionObserver | null = null;
const sortMode = ref<'recent' | 'popular'>('recent');

const alertDialog = ref<DialogState>({ visible: false, message: '' });
const photoPasswordDialog = ref<DialogState>({ visible: false, message: '' });
const pendingPhotoDelete = ref<PhotoUpload | null>(null);

const {
  heartCounts,
  initHearts,
  sendHeart,
  onContentUpdate,
  offContentUpdate,
} = useHearts();
const { award } = useAchievements();

const openKebab = ref<number | null>(null);

function toggleKebab(id: number) {
  openKebab.value = openKebab.value === id ? null : id;
}

function closeKebab() {
  openKebab.value = null;
}

const heartOverlayRef = ref<HTMLDivElement | null>(null);
let heartGame: Phaser.Game | null = null;

function showAlert(msg: string) {
  alertDialog.value = { visible: true, message: msg };
}

function onFileChange(e: Event) {
  const input = e.target as HTMLInputElement;
  for (const file of input.files ?? []) {
    queue.value.push({
      file,
      original: file,
      preview: URL.createObjectURL(file),
      status: 'pending',
    });
  }
  input.value = '';
}

function openEditor(index: number) {
  editingIndex.value = index;
  editorVisible.value = true;
}

function onEditorConfirm(editedFile: File) {
  const item =
    editingIndex.value !== null ? queue.value[editingIndex.value] : null;
  if (!item) return;
  URL.revokeObjectURL(item.preview);
  item.file = editedFile;
  item.preview = URL.createObjectURL(editedFile);
}

function removeItem(index: number) {
  const [removed] = queue.value.splice(index, 1);
  if (removed) URL.revokeObjectURL(removed.preview);
}

function clearQueue() {
  queue.value.forEach((item) => URL.revokeObjectURL(item.preview));
  queue.value = [];
}

async function submitPhoto() {
  if (!name.value.trim() || queue.value.length === 0 || uploading.value) return;

  if (uploadPassword.value.length < 5 || uploadPassword.value.length > 30) {
    showAlert(t('common.passwordLength'));
    return;
  }

  uploading.value = true;
  let lastError: unknown = null;
  await runWithConcurrency(queue.value, UPLOAD_CONCURRENCY, async (item) => {
    item.status = 'uploading';
    try {
      const wasEdited = item.file !== item.original;
      const image = await compressImage(item.file);
      const original = wasEdited ? await compressImage(item.original) : undefined;
      await uploadPhoto(name.value, image, uploadPassword.value, original);
      item.status = 'done';
    } catch (e) {
      item.status = 'error';
      lastError = e;
    }
  });

  const failed = queue.value.filter((item) => item.status === 'error');
  const succeeded = queue.value.length - failed.length;
  uploading.value = false;

  if (succeeded > 0) {
    award('photo');
    loadPhotos();
  }

  if (failed.length === 0) {
    clearQueue();
    name.value = '';
    uploadPassword.value = '';
    showAlert(t('photo.uploadSuccess'));
    return;
  }

  queue.value
    .filter((item) => item.status === 'done')
    .forEach((item) => URL.revokeObjectURL(item.preview));
  queue.value = failed;

  if (lastError instanceof ApiError && lastError.status === 503) {
    showAlert(t('photo.uploadDisabled'));
  } else if (succeeded === 0 && failed.length === 1) {
    showAlert(apiErrorMessage(lastError, t('photo.uploadFailed')));
  } else {
    showAlert(
      t('photo.uploadPartialFailed', { n: failed.length }, failed.length),
    );
  }
}

function startPhotoDelete(photo: PhotoUpload) {
  pendingPhotoDelete.value = photo;
  photoPasswordDialog.value = {
    visible: true,
    message: t('common.deletePrompt'),
  };
}

async function onPhotoDeleteConfirm(pw: string) {
  photoPasswordDialog.value.visible = false;
  if (!pw || !pendingPhotoDelete.value) return;

  try {
    await deleteUserPhoto(pendingPhotoDelete.value.id, pw);
    photos.value = photos.value.filter(
      (p) => p.id !== pendingPhotoDelete.value!.id,
    );
    showAlert(t('photo.deleted'));
  } catch (e) {
    showAlert(
      e instanceof ApiError && e.code === 'wrong_password'
        ? t('common.wrongPassword')
        : t('common.deleteFailed'),
    );
  }
  pendingPhotoDelete.value = null;
}

const lightboxIndex = ref<number | null>(null);

const lightboxItems = computed(() =>
  photos.value.map((p) => ({
    thumbnail: p.thumbnail,
    src: p.url,
    photoId: p.id,
  })),
);

function onLightboxHeart(photoId: number, event: MouseEvent) {
  sendHeart(photoId);
  award('heart');
  emitHeartAnimation(event);
}

async function loadPhotos() {
  try {
    const data = await getPhotos(sortMode.value);
    photos.value = data.items;
    nextOffset.value = data.next_offset;
    initHearts(photos.value);
  } catch {
    console.error('Failed to load photos');
  } finally {
    loadingPhotos.value = false;
  }
}

async function loadMorePhotos() {
  if (loadingMore.value || nextOffset.value === null) return;
  loadingMore.value = true;
  try {
    const data = await getPhotos(sortMode.value, nextOffset.value);
    photos.value.push(...data.items);
    nextOffset.value = data.next_offset;
    initHearts(data.items);
  } catch {
    console.error('Failed to load more photos');
  } finally {
    loadingMore.value = false;
  }
}

function switchSort(mode: 'recent' | 'popular') {
  if (sortMode.value === mode) return;
  sortMode.value = mode;
  loadingPhotos.value = true;
  photos.value = [];
  nextOffset.value = null;
  loadPhotos();
}

function onHeartClick(photo: PhotoUpload, event: MouseEvent | TouchEvent) {
  sendHeart(photo.id);
  award('heart');
  emitHeartAnimation(event);
}

function emitHeartAnimation(event: MouseEvent | TouchEvent) {
  const scene = heartGame?.scene.getScene('HeartScene') as
    | HeartScene
    | undefined;
  if (!scene) return;

  let x: number, y: number;
  if ('touches' in event && event.touches.length > 0) {
    x = event.touches[0].clientX;
    y = event.touches[0].clientY;
  } else {
    x = (event as MouseEvent).clientX;
    y = (event as MouseEvent).clientY;
  }

  const count = 1 + Math.floor(Math.random() * 3);
  for (let i = 0; i < count; i++) {
    setTimeout(() => scene.spawnHeart(x, y), i * 50);
  }
}

function initHeartOverlay() {
  if (!heartOverlayRef.value) return;
  heartGame = new Phaser.Game({
    type: Phaser.AUTO,
    parent: heartOverlayRef.value,
    width: window.innerWidth,
    height: window.innerHeight,
    transparent: true,
    scene: [HeartScene],
    input: false,
    audio: { noAudio: true },
    banner: false,
    scale: {
      mode: Phaser.Scale.RESIZE,
      autoCenter: Phaser.Scale.NO_CENTER,
    },
  });
}

const isDirty = computed(
  () => queue.value.length > 0 || editorVisible.value,
);

function onBeforeUnload(e: BeforeUnloadEvent) {
  if (isDirty.value) e.preventDefault();
}

onBeforeRouteLeave(() => {
  if (isDirty.value) {
    return window.confirm(t('photo.unsavedConfirm'));
  }
});

async function checkStorageStatus() {
  checkingStatus.value = true;
  uploadAvailable.value = await getPhotoStorageStatus();
  checkingStatus.value = false;
}

function setupPhotoScrollObserver() {
  photoScrollObserver?.disconnect();
  if (!photoSentinelRef.value) return;
  photoScrollObserver = new IntersectionObserver(
    (entries) => {
      if (entries[0].isIntersecting) loadMorePhotos();
    },
    { rootMargin: '200px' },
  );
  photoScrollObserver.observe(photoSentinelRef.value);
}

watch(photoSentinelRef, () => setupPhotoScrollObserver());

onMounted(() => {
  window.addEventListener('beforeunload', onBeforeUnload);
  document.addEventListener('click', closeKebab);
  if (timeGateOpen.value) {
    checkStorageStatus();
  } else {
    updateCountdown();
    countdownTimer = setInterval(updateCountdown, 1000);
  }
  loadPhotos();
  onContentUpdate('photos', loadPhotos);
  initHeartOverlay();
});

onUnmounted(() => {
  window.removeEventListener('beforeunload', onBeforeUnload);
  document.removeEventListener('click', closeKebab);
  if (countdownTimer) clearInterval(countdownTimer);
  offContentUpdate('photos', loadPhotos);
  photoScrollObserver?.disconnect();
  heartGame?.destroy(true);
  heartGame = null;
});
</script>

<style scoped>
.sort-btn {
  border: 3px solid var(--btn-border);
  cursor: pointer;
  transition: background-color 0.15s ease;
}

.sort-btn-active {
  background-color: var(--btn-bg);
  color: white;
  box-shadow:
    inset 2px 2px 0 var(--btn-highlight),
    inset -2px -2px 0 var(--btn-shadow);
}

.sort-btn-inactive {
  background-color: var(--btn-bg);
  opacity: 0.45;
  color: white;
}

.sort-btn-inactive:hover {
  opacity: 0.65;
}

.masonry-grid {
  display: flex;
  gap: 0.5rem;
}

.masonry-col {
  flex: 1;
  min-width: 0;
}

.heart-btn {
  transition: transform 0.15s ease;
}

.heart-btn:active {
  transform: scale(1.3);
}

.heart-icon {
  color: #ff4466;
  font-size: 18px;
  transition: transform 0.15s ease;
}

.heart-btn:active .heart-icon {
  transform: scale(1.4);
}
</style>
