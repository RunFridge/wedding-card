<template>
  <Teleport to="body">
    <Transition name="lightbox-fade">
      <div
        v-if="modelValue !== null"
        role="dialog"
        aria-modal="true"
        :aria-label="t('lightbox.viewer')"
        class="fixed inset-0 z-[100] flex items-center justify-center bg-black/90 touch-none"
        @click.self="close"
        @touchstart.passive="onTouchStart"
        @touchmove.passive="onTouchMove"
        @touchend.passive="onTouchEnd"
      >
        <!-- Sliding image -->
        <Transition :name="slideDirection">
          <img
            :key="modelValue"
            :src="currentItem?.src"
            :alt="
              t('lightbox.photoAlt', {
                current: (modelValue ?? 0) + 1,
                total: items.length,
              })
            "
            class="absolute inset-0 w-full h-full object-contain"
          />
        </Transition>

        <!-- Close button -->
        <button
          :aria-label="t('common.close')"
          @click="close"
          class="absolute top-4 right-4 w-10 h-10 bg-wood-dark/80 text-parchment rounded-full flex items-center justify-center text-xl leading-none cursor-pointer border-2 border-parchment/50 z-10"
        >
          &times;
        </button>

        <!-- Left arrow -->
        <button
          v-if="items.length > 1"
          :aria-label="t('lightbox.prev')"
          @click.stop="navigate(-1)"
          class="absolute left-3 top-1/2 -translate-y-1/2 w-10 h-10 bg-wood-dark/70 text-parchment rounded-full flex items-center justify-center text-xl cursor-pointer border-2 border-parchment/40 z-10 hover:bg-wood-dark/90 transition-colors"
        >
          &#8249;
        </button>

        <!-- Right arrow -->
        <button
          v-if="items.length > 1"
          :aria-label="t('lightbox.next')"
          @click.stop="navigate(1)"
          class="absolute right-3 top-1/2 -translate-y-1/2 w-10 h-10 bg-wood-dark/70 text-parchment rounded-full flex items-center justify-center text-xl cursor-pointer border-2 border-parchment/40 z-10 hover:bg-wood-dark/90 transition-colors"
        >
          &#8250;
        </button>

        <!-- Bottom bar -->
        <div
          class="absolute bottom-4 left-1/2 -translate-x-1/2 flex items-center gap-3 z-10"
        >
          <button
            v-if="heartCounts && currentPhotoId !== null"
            :aria-label="t('common.like')"
            class="lightbox-heart-btn flex items-center gap-1 bg-wood-dark/70 text-parchment rounded-full px-3 py-1.5 border border-parchment/40"
            @click.stop="$emit('heart', currentPhotoId!, $event)"
          >
            <span class="lightbox-heart-icon">&#9829;</span>
            <span class="text-sm">{{ heartCounts[currentPhotoId!] || 0 }}</span>
          </button>
          <span class="text-parchment/70 text-sm"
            >{{ modelValue + 1 }} / {{ items.length }}</span
          >
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

export interface LightboxItem {
  thumbnail?: string;
  src: string;
  photoId?: number;
}

const props = defineProps<{
  items: LightboxItem[];
  modelValue: number | null;
  heartCounts?: Record<number, number>;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: number | null];
  heart: [photoId: number, event: MouseEvent];
}>();

const slideDirection = ref('slide-left');

const currentItem = computed(() =>
  props.modelValue !== null ? props.items[props.modelValue] : null,
);

const currentPhotoId = computed(() => currentItem.value?.photoId ?? null);

function preloadNearby() {
  if (props.modelValue === null || props.items.length <= 1) return;
  for (let offset = -3; offset <= 3; offset++) {
    if (offset === 0) continue;
    const idx =
      (props.modelValue + offset + props.items.length) % props.items.length;
    const src = props.items[idx]?.src;
    if (src) {
      const img = new Image();
      img.src = src;
    }
  }
}

watch(
  () => props.modelValue,
  () => {
    preloadNearby();
  },
);

function close() {
  emit('update:modelValue', null);
}

function navigate(direction: number) {
  if (props.modelValue === null) return;
  const next =
    (props.modelValue + direction + props.items.length) % props.items.length;
  if (next === props.modelValue) return;
  slideDirection.value = direction > 0 ? 'slide-left' : 'slide-right';
  emit('update:modelValue', next);
}

function handleKeydown(e: KeyboardEvent) {
  if (props.modelValue === null) return;
  if (e.key === 'ArrowLeft') navigate(-1);
  else if (e.key === 'ArrowRight') navigate(1);
  else if (e.key === 'Escape') close();
}

let touchStartX = 0;
let touchStartY = 0;
let swiping = false;

function onTouchStart(e: TouchEvent) {
  touchStartX = e.touches[0].clientX;
  touchStartY = e.touches[0].clientY;
  swiping = true;
}

function onTouchMove(e: TouchEvent) {
  if (!swiping) return;
  const dy = Math.abs(e.touches[0].clientY - touchStartY);
  if (dy > 60) swiping = false;
}

function onTouchEnd(e: TouchEvent) {
  if (!swiping) return;
  swiping = false;
  const dx = e.changedTouches[0].clientX - touchStartX;
  if (Math.abs(dx) < 50) return;
  const dir = dx < 0 ? 1 : -1;
  slideDirection.value = dir > 0 ? 'slide-left' : 'slide-right';
  navigate(dir);
}

onMounted(() => window.addEventListener('keydown', handleKeydown));
onUnmounted(() => window.removeEventListener('keydown', handleKeydown));
</script>

<style scoped>
.lightbox-fade-enter-active,
.lightbox-fade-leave-active {
  transition: opacity 0.2s ease;
}

.lightbox-fade-enter-from,
.lightbox-fade-leave-to {
  opacity: 0;
}

.slide-left-enter-active,
.slide-left-leave-active,
.slide-right-enter-active,
.slide-right-leave-active {
  transition:
    transform 0.25s ease,
    opacity 0.25s ease;
}

.slide-left-enter-from {
  transform: translateX(60px);
  opacity: 0;
}

.slide-left-leave-to {
  transform: translateX(-60px);
  opacity: 0;
}

.slide-right-enter-from {
  transform: translateX(-60px);
  opacity: 0;
}

.slide-right-leave-to {
  transform: translateX(60px);
  opacity: 0;
}

.lightbox-heart-btn {
  transition: transform 0.15s ease;
}

.lightbox-heart-btn:active {
  transform: scale(1.3);
}

.lightbox-heart-icon {
  color: #ff4466;
  font-size: 14px;
  transition: transform 0.15s ease;
}

.lightbox-heart-btn:active .lightbox-heart-icon {
  transform: scale(1.4);
}
</style>
