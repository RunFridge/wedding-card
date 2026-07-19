<template>
  <div class="max-w-md mx-auto px-4 py-6">
    <div class="wooden-panel p-2 mb-6">
      <div class="wooden-panel-inner px-4 py-3 text-center">
        <h1 class="text-xl text-secondary mb-1">{{ t('map.title') }}</h1>
        <p class="text-wood-dark text-sm">
          {{ t('map.subtitle') }}
        </p>
      </div>
    </div>

    <div v-if="hasEmbed" class="wooden-panel p-2 mb-6">
      <div class="wooden-panel-inner p-2">
        <MapEmbed
          :provider="MAP_PROVIDERS.embed_provider"
          :latitude="MAP_PROVIDERS.latitude"
          :longitude="MAP_PROVIDERS.longitude"
          :api-key="MAP_PROVIDERS.api_key"
        />
      </div>
    </div>

    <div class="wooden-panel p-2 mb-6">
      <div class="wooden-panel-inner p-4">
        <h2 class="text-secondary text-base mb-2">{{ t('map.address') }}</h2>
        <p class="text-secondary text-lg font-semibold">
          {{ VENUE_NAME }} {{ VENUE_FLOOR }} {{ VENUE_HALL }}
        </p>
        <p class="text-wood-dark/80 text-lg mb-3">{{ VENUE_ADDRESS }}</p>

        <button
          @click="copyAddress"
          class="pixel-btn-outline bg-parchment w-full text-sm mb-3"
        >
          {{ copied ? t('map.copied') : t('map.copyAddress') }}
        </button>

        <div
          v-if="linkEntries.length"
          class="flex justify-center items-center gap-4"
        >
          <a
            v-for="entry in linkEntries"
            :key="entry.key"
            :href="entry.link"
            target="_blank"
            rel="noopener noreferrer"
            :aria-label="
              t('map.openInAria', { provider: t(`map.providers.${entry.key}`) })
            "
          >
            <img
              :src="providerIcons[entry.key]"
              :alt="t(`map.providers.${entry.key}`)"
              class="h-12 rounded-lg cursor-pointer"
            />
          </a>
        </div>
      </div>
    </div>

    <div v-if="CHARTER_BUS.length" class="wooden-panel p-2 mb-6">
      <div class="wooden-panel-inner p-4 text-center">
        <h2 class="text-secondary text-base mb-3">
          <TwEmoji emoji="🚍" size="1rem" /> {{ t('map.charterBus') }}
        </h2>
        <div
          v-for="(bus, i) in CHARTER_BUS"
          :key="i"
          :class="i < CHARTER_BUS.length - 1 ? 'mb-3' : ''"
        >
          <p class="text-secondary text-lg font-semibold whitespace-pre-line">
            {{ bus.location.replace(/ \(/g, '\n(') }}
          </p>
          <p class="text-wood-dark text-sm">
            {{ bus.company }} {{ bus.bus_number }}
          </p>
          <p class="text-wood-dark text-base font-semibold">
            {{ bus.departure }}
          </p>
        </div>

        <div
          v-if="CHARTER_BUS_NOTICE"
          class="parchment-bg bus-notice p-3 mt-4 text-left"
        >
          <p class="text-secondary text-sm font-bold mb-2 text-center">
            <TwEmoji emoji="📢" size="1rem" /> {{ t('map.charterBusNotice') }}
          </p>
          <p
            class="text-wood-dark text-xs break-keep leading-relaxed whitespace-pre-line"
          >
            {{ CHARTER_BUS_NOTICE }}
          </p>
        </div>
      </div>
    </div>

    <div class="wooden-panel p-2 mb-6">
      <div class="wooden-panel-inner p-4">
        <h2 class="text-secondary text-sm mb-3">{{ t('map.transport') }}</h2>

        <div class="space-y-3">
          <div class="parchment-bg p-3">
            <h3 class="text-secondary text-xs mb-2">
              <TwEmoji emoji="🚇" size="0.875rem" /> {{ t('map.subway') }}
            </h3>
            <p
              v-for="line in SUBWAY_INFO"
              :key="line"
              class="text-wood-dark/80 text-sm"
            >
              {{ line }}
            </p>
          </div>

          <div class="parchment-bg p-3">
            <h3 class="text-secondary text-xs mb-2">
              <TwEmoji emoji="🚌" size="0.875rem" /> {{ t('map.bus') }}
            </h3>
            <div
              v-for="(bus, i) in BUS_INFO"
              :key="bus.stop"
              :class="i < BUS_INFO.length - 1 ? 'mb-2' : ''"
            >
              <p class="text-wood-dark text-sm font-semibold">{{ bus.stop }}</p>
              <p class="text-wood-dark/80 text-sm">{{ bus.routes }}</p>
            </div>
          </div>

          <div class="parchment-bg p-3">
            <h3 class="text-secondary text-xs mb-2">
              <TwEmoji emoji="🚗" size="0.875rem" /> {{ t('map.car') }}
            </h3>
            <p class="text-wood-dark/80 text-sm">{{ CAR_INFO }}</p>
          </div>
        </div>
      </div>
    </div>

    <div class="wooden-panel p-2">
      <div class="wooden-panel-inner p-4">
        <h2 class="text-secondary text-sm mb-2">{{ t('map.venueContact') }}</h2>
        <p class="text-wood-dark/80 text-sm mb-1">{{ VENUE_NAME }}</p>
        <a
          :href="`tel:${VENUE_PHONE.replace(/-/g, '')}`"
          class="text-wood-dark text-xs underline"
          >{{ VENUE_PHONE }}</a
        >
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import TwEmoji from '../components/TwEmoji.vue';
import MapEmbed from '../components/MapEmbed.vue';
import tmapIcon from '../assets/icons/tmap.png';
import kakaomapIcon from '../assets/icons/kakaomap.png';
import navermapIcon from '../assets/icons/navermap.png';
import googlemapIcon from '../assets/icons/googlemap.png';
import {
  VENUE_NAME,
  VENUE_ADDRESS,
  VENUE_FLOOR,
  VENUE_HALL,
  VENUE_PHONE,
  MAP_PROVIDERS,
  SUBWAY_INFO,
  BUS_INFO,
  CAR_INFO,
  CHARTER_BUS,
  CHARTER_BUS_NOTICE,
} from '../config/wedding';

const { t } = useI18n();

const providerIcons: Record<string, string> = {
  tmap: tmapIcon,
  google: googlemapIcon,
  kakao: kakaomapIcon,
  naver: navermapIcon,
};

const hasEmbed = computed(
  () => MAP_PROVIDERS.embed_provider && MAP_PROVIDERS.api_key,
);

const linkEntries = computed(() =>
  (['tmap', 'google', 'kakao', 'naver'] as const)
    .filter((key) => MAP_PROVIDERS.links[key])
    .map((key) => ({ key, link: MAP_PROVIDERS.links[key] })),
);

const copied = ref(false);

function copyAddress() {
  navigator.clipboard.writeText(VENUE_ADDRESS);
  copied.value = true;
  setTimeout(() => {
    copied.value = false;
  }, 2000);
}
</script>

<style scoped>
.parchment-bg.bus-notice {
  border-left-width: 4px;
  border-left-color: rgb(var(--c-primary));
  box-shadow:
    0 0 8px rgba(139, 105, 20, 0.18),
    inset 0 0 10px rgba(139, 105, 20, 0.08);
}

html.dark .parchment-bg.bus-notice {
  box-shadow:
    0 0 10px rgba(196, 148, 58, 0.22),
    inset 0 0 10px rgba(0, 0, 0, 0.3);
}
</style>
