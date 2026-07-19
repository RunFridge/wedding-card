<template>
  <Teleport to="body">
    <Transition name="dialog">
      <div
        v-if="visible"
        role="dialog"
        aria-modal="true"
        aria-labelledby="stardew-dialog-msg"
        class="fixed inset-0 z-[100] flex items-center justify-center px-4"
        @click.self="cancel"
      >
        <div class="fixed inset-0 bg-black/50"></div>

        <div class="relative wooden-panel p-3 w-full max-w-xs">
          <div class="wooden-panel-inner px-5 py-5 text-center">
            <p
              id="stardew-dialog-msg"
              class="text-secondary text-sm leading-relaxed mb-4 whitespace-pre-wrap break-keep"
            >
              {{ message }}
            </p>

            <input
              v-if="inputMode === 'password'"
              ref="inputRef"
              v-model="inputValue"
              type="password"
              :aria-label="t('dialog.password')"
              class="pixel-input w-full mb-4"
              :placeholder="t('dialog.passwordPlaceholder')"
              @keyup.enter="confirm"
            />

            <input
              v-if="inputMode === 'name'"
              ref="inputRef"
              v-model="inputValue"
              type="text"
              maxlength="30"
              :aria-label="t('dialog.name')"
              class="pixel-input w-full mb-4"
              :placeholder="t('common.namePlaceholder')"
              @keyup.enter="confirm"
            />

            <textarea
              v-if="inputMode === 'edit'"
              ref="inputRef"
              v-model="inputValue"
              rows="4"
              maxlength="500"
              :aria-label="t('dialog.message')"
              class="pixel-input w-full resize-none mb-1"
              :placeholder="t('dialog.messagePlaceholder')"
            ></textarea>
            <p
              v-if="inputMode === 'edit'"
              class="text-xs text-wood-dark/60 text-right mb-3"
            >
              {{ inputValue.length }}/500
            </p>

            <div v-if="inputMode" class="flex gap-2">
              <button
                class="pixel-btn text-sm px-4 flex-1 !bg-wood-dark/20"
                @click="cancel"
              >
                {{ t('common.cancel') }}
              </button>
              <button class="pixel-btn text-sm px-4 flex-1" @click="confirm">
                {{ t('common.confirm') }}
              </button>
            </div>
            <button v-else class="pixel-btn text-sm px-6" @click="close">
              {{ t('common.confirm') }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

type InputMode = 'password' | 'edit' | 'name' | 'confirm' | null;

const props = withDefaults(
  defineProps<{
    visible?: boolean;
    message?: string;
    inputMode?: InputMode;
    initialValue?: string;
  }>(),
  {
    visible: false,
    message: '',
    inputMode: null,
    initialValue: '',
  },
);

const emit = defineEmits<{
  close: [];
  cancel: [];
  confirm: [value: string];
}>();

const inputValue = ref('');
const inputRef = ref<HTMLInputElement | HTMLTextAreaElement | null>(null);

watch(
  () => props.visible,
  (val) => {
    if (val) {
      inputValue.value = props.initialValue || '';
      nextTick(() => inputRef.value?.focus());
    }
  },
);

function close() {
  emit('close');
}

function cancel() {
  emit('cancel');
  emit('close');
}

function confirm() {
  emit('confirm', inputValue.value);
}
</script>

<style scoped>
.dialog-enter-active,
.dialog-leave-active {
  transition: opacity 0.2s ease;
}

.dialog-enter-from,
.dialog-leave-to {
  opacity: 0;
}
</style>
