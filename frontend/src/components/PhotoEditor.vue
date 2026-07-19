<template>
  <Teleport to="body">
    <Transition name="editor-fade">
      <div v-if="visible" class="fixed inset-0 z-[100] flex flex-col bg-black">
        <!-- Top bar -->
        <div
          class="flex items-center justify-between px-2 py-2 bg-wood-dark border-b-2 border-wood/50 shrink-0"
        >
          <button class="pixel-btn text-xs px-3 py-1" @click="cancel">
            {{ t('common.cancel') }}
          </button>
          <div class="flex items-center gap-1">
            <button
              class="pixel-btn text-xs px-2 py-1"
              :class="{ 'opacity-30 pointer-events-none': !canUndo }"
              @click="undo"
            >
              ↩
            </button>
            <button
              class="pixel-btn text-xs px-2 py-1"
              :class="{ 'opacity-30 pointer-events-none': !canRedo }"
              @click="redo"
            >
              ↪
            </button>
            <button class="pixel-btn text-xs px-2 py-1" @click="resetToDefault">
              {{ t('photoEditor.reset') }}
            </button>
          </div>
          <button class="pixel-btn text-xs px-3 py-1" @click="confirm">
            {{ t('photoEditor.done') }}
          </button>
        </div>

        <!-- Canvas area -->
        <div class="flex-1 min-h-0 relative">
          <div ref="cropperContainer" class="absolute inset-0 overflow-hidden">
            <img
              ref="sourceImage"
              :src="sourceUrl"
              class="block max-w-full"
              @load="onImageLoad"
            />
          </div>
        </div>

        <!-- Bottom toolbar -->
        <div class="bg-wood-dark border-t-2 border-wood/50 shrink-0">
          <!-- Tab switcher -->
          <div class="flex border-b border-wood/30">
            <button
              v-for="tab in tabs"
              :key="tab.id"
              class="flex-1 py-2 text-xs text-center transition-colors"
              :class="
                activeTab === tab.id
                  ? 'text-parchment bg-wood/30'
                  : 'text-parchment/50'
              "
              @click="activeTab = tab.id"
            >
              {{ t(`photoEditor.tabs.${tab.id}`) }}
            </button>
          </div>

          <!-- Crop tab -->
          <div v-if="activeTab === 'crop'" class="px-3 py-3">
            <div class="flex items-center gap-2 justify-center mb-2">
              <button
                v-for="ratio in aspectRatios"
                :key="ratio.id"
                class="pixel-btn text-xs px-2 py-1"
                :class="{ '!bg-wood': selectedRatio === ratio.id }"
                @click="setAspectRatio(ratio)"
              >
                {{ t(`photoEditor.ratios.${ratio.id}`) }}
              </button>
            </div>
            <div class="flex items-center gap-3 justify-center">
              <button class="pixel-btn text-xs px-3 py-1" @click="rotate(-90)">
                ↺ 90°
              </button>
              <button class="pixel-btn text-xs px-3 py-1" @click="rotate(90)">
                ↻ 90°
              </button>
            </div>
          </div>

          <!-- Adjust tab -->
          <div v-if="activeTab === 'adjust'" class="px-3 py-3">
            <!-- Filter presets -->
            <div class="flex gap-2 overflow-x-auto pb-2 mb-2 -mx-1 px-1">
              <button
                v-for="preset in filterPresets"
                :key="preset.name"
                class="flex-shrink-0 px-3 py-1 text-xs rounded-full border transition-colors"
                :class="
                  activePreset === preset.name
                    ? 'bg-wood text-parchment border-wood'
                    : 'bg-transparent text-parchment/70 border-parchment/30'
                "
                @click="applyPreset(preset)"
              >
                {{ t(`photoEditor.filters.${preset.name}`) }}
              </button>
            </div>
            <!-- Sliders -->
            <div
              v-for="slider in sliders"
              :key="slider.key"
              class="flex items-center gap-2 mb-1.5"
            >
              <span class="text-parchment/70 text-xs w-12 text-right">{{
                t(`photoEditor.adjustments.${slider.key}`)
              }}</span>
              <input
                type="range"
                :min="slider.min"
                :max="slider.max"
                :step="slider.step"
                :value="adjustments[slider.key]"
                class="slider flex-1"
                @input="onSliderInput(slider.key, $event)"
                @change="pushSnapshot"
              />
              <span class="text-parchment/50 text-xs w-8"
                >{{ Math.round(adjustments[slider.key] * 100) }}%</span
              >
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, onUnmounted } from 'vue';
import { useI18n } from 'vue-i18n';
import Cropper from 'cropperjs';
import 'cropperjs/dist/cropper.css';
import type { ImageAdjustments, FilterPreset } from '../types';

const { t } = useI18n();

interface EditorSnapshot {
  adjustments: ImageAdjustments;
  activePreset: string;
  rotationDeg: number;
}

const props = defineProps<{
  visible: boolean;
  sourceUrl: string;
}>();

const emit = defineEmits<{
  'update:visible': [value: boolean];
  confirm: [editedFile: File];
}>();

const sourceImage = ref<HTMLImageElement | null>(null);
const cropperContainer = ref<HTMLDivElement | null>(null);

let cropper: Cropper | null = null;

const activeTab = ref<'crop' | 'adjust'>('crop');
const currentRotation = ref(0);

const tabs = [{ id: 'crop' as const }, { id: 'adjust' as const }];

const aspectRatios = [
  { id: 'free', value: NaN },
  { id: 'square', value: 1 },
  { id: 'landscape', value: 4 / 3 },
  { id: 'portrait', value: 3 / 4 },
];

const selectedRatio = ref('free');

const filterPresets: FilterPreset[] = [
  { name: 'original', filter: '' },
  {
    name: 'warm',
    filter: 'brightness(1.05) saturate(1.3) sepia(0.15)',
  },
  {
    name: 'cool',
    filter: 'brightness(1.05) saturate(0.9) hue-rotate(15deg)',
  },
  { name: 'sepia', filter: 'sepia(0.8)' },
  { name: 'grayscale', filter: 'grayscale(1)' },
  {
    name: 'vivid',
    filter: 'brightness(1.1) contrast(1.2) saturate(1.4)',
  },
];

const activePreset = ref('original');

const adjustments = ref<ImageAdjustments>({
  brightness: 1,
  contrast: 1,
  saturation: 1,
});

const sliders = [
  {
    key: 'brightness' as keyof ImageAdjustments,
    min: 0.5,
    max: 1.5,
    step: 0.01,
  },
  {
    key: 'contrast' as keyof ImageAdjustments,
    min: 0.5,
    max: 1.5,
    step: 0.01,
  },
  {
    key: 'saturation' as keyof ImageAdjustments,
    min: 0,
    max: 2,
    step: 0.01,
  },
];

// --- History (undo / redo) ---

const history = ref<EditorSnapshot[]>([]);
const historyIndex = ref(-1);

const canUndo = computed(() => historyIndex.value > 0);
const canRedo = computed(() => historyIndex.value < history.value.length - 1);

function captureSnapshot(): EditorSnapshot {
  return {
    adjustments: { ...adjustments.value },
    activePreset: activePreset.value,
    rotationDeg: currentRotation.value,
  };
}

function pushSnapshot() {
  history.value = history.value.slice(0, historyIndex.value + 1);
  history.value.push(captureSnapshot());
  historyIndex.value = history.value.length - 1;
}

function applySnapshot(snapshot: EditorSnapshot) {
  const rotationDelta = snapshot.rotationDeg - currentRotation.value;
  if (rotationDelta !== 0 && cropper) {
    cropper.rotateTo(snapshot.rotationDeg);
  }
  currentRotation.value = snapshot.rotationDeg;
  adjustments.value = { ...snapshot.adjustments };
  activePreset.value = snapshot.activePreset;
}

function undo() {
  if (!canUndo.value) return;
  historyIndex.value--;
  applySnapshot(history.value[historyIndex.value]);
}

function redo() {
  if (!canRedo.value) return;
  historyIndex.value++;
  applySnapshot(history.value[historyIndex.value]);
}

// --- Filter ---

const filterString = computed(() => {
  const { brightness, contrast, saturation } = adjustments.value;
  const base = `brightness(${brightness}) contrast(${contrast}) saturate(${saturation})`;
  const preset = filterPresets.find((p) => p.name === activePreset.value);
  if (!preset || !preset.filter) return base;
  return `${base} ${preset.filter}`;
});

function applyFilterToDOM() {
  if (!cropperContainer.value) return;
  const imgs = cropperContainer.value.querySelectorAll<HTMLImageElement>(
    '.cropper-canvas img, .cropper-view-box img',
  );
  imgs.forEach((img) => {
    img.style.filter = filterString.value;
  });
}

watch(filterString, applyFilterToDOM);

// --- Actions ---

function onSliderInput(key: keyof ImageAdjustments, e: Event) {
  const value = parseFloat((e.target as HTMLInputElement).value);
  adjustments.value[key] = value;
  activePreset.value = 'original';
}

function applyPreset(preset: FilterPreset) {
  activePreset.value = preset.name;
  if (preset.name === 'original') {
    adjustments.value = { brightness: 1, contrast: 1, saturation: 1 };
  }
  pushSnapshot();
}

function setAspectRatio(ratio: { id: string; value: number }) {
  selectedRatio.value = ratio.id;
  cropper?.setAspectRatio(ratio.value);
}

function rotate(degrees: number) {
  cropper?.rotate(degrees);
  currentRotation.value += degrees;
  pushSnapshot();
}

function resetToDefault() {
  adjustments.value = { brightness: 1, contrast: 1, saturation: 1 };
  activePreset.value = 'original';
  selectedRatio.value = 'free';
  currentRotation.value = 0;
  cropper?.reset();
  cropper?.setAspectRatio(NaN);

  history.value = [];
  historyIndex.value = -1;
  pushSnapshot();
}

// --- Cropper lifecycle ---

function onImageLoad() {
  if (cropper) return;
  initCropper();
}

function initCropper() {
  if (cropper || !sourceImage.value) return;
  cropper = new Cropper(sourceImage.value, {
    viewMode: 1,
    autoCropArea: 1,
    responsive: true,
    background: false,
    checkCrossOrigin: false,
    ready: applyFilterToDOM,
  });

  if (history.value.length === 0) {
    pushSnapshot();
  }
}

function destroyCropper() {
  if (cropper) {
    cropper.destroy();
    cropper = null;
  }
}

function resetState() {
  activeTab.value = 'crop';
  activePreset.value = 'original';
  selectedRatio.value = 'free';
  currentRotation.value = 0;
  adjustments.value = { brightness: 1, contrast: 1, saturation: 1 };
  history.value = [];
  historyIndex.value = -1;
}

watch(
  () => props.visible,
  (val) => {
    if (val) {
      resetState();
      nextTick(() => {
        if (sourceImage.value?.complete && sourceImage.value.naturalWidth > 0) {
          initCropper();
        }
      });
    } else {
      destroyCropper();
    }
  },
);

onUnmounted(destroyCropper);

function cancel() {
  destroyCropper();
  emit('update:visible', false);
}

async function confirm() {
  if (!cropper) return;

  try {
    const croppedCanvas = cropper.getCroppedCanvas({
      imageSmoothingEnabled: true,
      imageSmoothingQuality: 'high',
    });

    const needsFilter =
      filterString.value !== 'brightness(1) contrast(1) saturate(1)';

    let outputCanvas: HTMLCanvasElement;
    if (needsFilter) {
      outputCanvas = document.createElement('canvas');
      outputCanvas.width = croppedCanvas.width;
      outputCanvas.height = croppedCanvas.height;
      const ctx = outputCanvas.getContext('2d')!;
      ctx.filter = filterString.value;
      ctx.drawImage(croppedCanvas, 0, 0);
    } else {
      outputCanvas = croppedCanvas;
    }

    const blob = await new Promise<Blob>((resolve, reject) => {
      outputCanvas.toBlob(
        (b) => (b ? resolve(b) : reject(new Error('Failed to export image'))),
        'image/jpeg',
        0.92,
      );
    });

    const file = new File([blob], 'edited-photo.jpg', { type: 'image/jpeg' });
    destroyCropper();
    emit('confirm', file);
    emit('update:visible', false);
  } catch (err) {
    console.error('Failed to export edited image:', err);
  }
}
</script>

<style scoped>
.editor-fade-enter-active,
.editor-fade-leave-active {
  transition: opacity 0.2s ease;
}

.editor-fade-enter-from,
.editor-fade-leave-to {
  opacity: 0;
}

.slider {
  -webkit-appearance: none;
  appearance: none;
  height: 4px;
  border-radius: 2px;
  background: rgba(245, 230, 200, 0.2);
  outline: none;
}

.slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #a0722a;
  border: 2px solid #f5e6c8;
  cursor: pointer;
}

.slider::-moz-range-thumb {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #a0722a;
  border: 2px solid #f5e6c8;
  cursor: pointer;
}
</style>
