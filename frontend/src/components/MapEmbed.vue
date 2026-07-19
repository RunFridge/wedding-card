<template>
  <div class="aspect-video bg-parchment-dark">
    <iframe
      v-if="provider === 'google'"
      :src="googleEmbedUrl"
      width="100%"
      height="100%"
      style="border: 0"
      allowfullscreen
      loading="lazy"
      referrerpolicy="no-referrer-when-downgrade"
    ></iframe>
    <div v-else ref="mapContainer" class="w-full h-full"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';

const props = defineProps<{
  provider: string;
  latitude: number;
  longitude: number;
  apiKey: string;
}>();

const mapContainer = ref<HTMLElement | null>(null);

const googleEmbedUrl = computed(
  () =>
    `https://www.google.com/maps/embed/v1/place?key=${props.apiKey}&q=${props.latitude},${props.longitude}&zoom=15`,
);

function loadScript(src: string): Promise<void> {
  return new Promise((resolve, reject) => {
    if (document.querySelector(`script[src="${src}"]`)) {
      resolve();
      return;
    }
    const script = document.createElement('script');
    script.src = src;
    script.onload = () => resolve();
    script.onerror = () => reject(new Error(`Failed to load ${src}`));
    document.head.appendChild(script);
  });
}

onMounted(async () => {
  if (props.provider === 'google' || !mapContainer.value) return;

  try {
    switch (props.provider) {
      case 'kakao': {
        await loadScript(
          `https://dapi.kakao.com/v2/maps/sdk.js?appkey=${props.apiKey}&autoload=false`,
        );
        const kakao = (window as any).kakao;
        kakao.maps.load(() => {
          const map = new kakao.maps.Map(mapContainer.value, {
            center: new kakao.maps.LatLng(props.latitude, props.longitude),
            level: 3,
          });
          new kakao.maps.Marker({
            map,
            position: new kakao.maps.LatLng(props.latitude, props.longitude),
          });
        });
        break;
      }
      case 'naver': {
        await loadScript(
          `https://oapi.map.naver.com/openapi/v3/maps.js?ncpKeyId=${props.apiKey}`,
        );
        const naver = (window as any).naver;
        const position = new naver.maps.LatLng(props.latitude, props.longitude);
        const map = new naver.maps.Map(mapContainer.value, {
          center: position,
          zoom: 16,
        });
        new naver.maps.Marker({ map, position });
        break;
      }
      case 'tmap': {
        await loadScript(
          `https://apis.openapi.sk.com/tmap/jsv2?version=1&appKey=${props.apiKey}`,
        );
        const Tmapv2 = (window as any).Tmapv2;
        const map = new Tmapv2.Map(mapContainer.value, {
          center: new Tmapv2.LatLng(props.latitude, props.longitude),
          zoom: 16,
          width: '100%',
          height: '100%',
        });
        new Tmapv2.Marker({
          map,
          position: new Tmapv2.LatLng(props.latitude, props.longitude),
        });
        break;
      }
    }
  } catch (e) {
    console.error('Failed to load map embed:', e);
  }
});
</script>
