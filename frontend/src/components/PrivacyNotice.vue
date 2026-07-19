<template>
  <button
    class="fixed left-1/2 -translate-x-1/2 z-[60] text-[9px] text-parchment/40 hover:text-parchment/70 cursor-pointer"
    style="bottom: calc(4.25rem + env(safe-area-inset-bottom, 0px))"
    @click="open = true"
  >
    {{ t('privacy.title') }}
  </button>

  <Teleport to="body">
    <Transition name="dialog">
      <div
        v-if="open"
        role="dialog"
        aria-modal="true"
        aria-labelledby="privacy-title"
        class="fixed inset-0 z-[300] flex items-center justify-center px-4"
        @click.self="open = false"
      >
        <div class="fixed inset-0 bg-black/50"></div>
        <div
          class="relative wooden-panel p-3 w-full max-w-sm max-h-[80vh] overflow-y-auto"
        >
          <div class="wooden-panel-inner px-5 py-5">
            <p
              id="privacy-title"
              class="text-primary text-base font-bold mb-3 text-center"
            >
              {{ t('privacy.title') }}
            </p>
            <div
              class="text-secondary text-xs leading-relaxed space-y-2 break-keep"
            >
              <p>{{ t('privacy.intro') }}</p>
              <p>
                <b>{{ t('privacy.collectLabel') }}</b
                >:
                {{ t('privacy.collectBody') }}
              </p>
              <p>
                <b>{{ t('privacy.purposeLabel') }}</b
                >:
                {{ t('privacy.purposeBody') }}
              </p>
              <p>
                <b>{{ t('privacy.retentionLabel') }}</b
                >:
                {{ t('privacy.retentionBody') }}
              </p>
              <p>{{ t('privacy.selfManage') }}</p>
            </div>
            <button
              class="pixel-btn text-sm px-6 w-full mt-4"
              @click="open = false"
            >
              {{ t('common.confirm') }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const open = ref(false);
</script>
