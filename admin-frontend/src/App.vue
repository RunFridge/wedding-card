<script setup lang="ts">
import { computed, watch, onUnmounted } from 'vue';
import { useRoute } from 'vue-router';
import { useAuth } from '@/lib/auth';
import { useAdminWS } from '@/composables/useAdminWS';
import AppSidebar from '@/components/AppSidebar.vue';
import AppHeader from '@/components/AppHeader.vue';

const route = useRoute();
const { authenticated } = useAuth();
const { connect, disconnect } = useAdminWS();

const showLayout = computed(
  () => authenticated.value && route.name !== 'login' && route.name !== 'setup',
);

watch(
  showLayout,
  (val) => {
    if (val) connect();
    else disconnect();
  },
  { immediate: true },
);

onUnmounted(disconnect);
</script>

<template>
  <div v-if="showLayout" class="flex h-screen">
    <AppSidebar />
    <div class="flex flex-1 flex-col overflow-hidden">
      <AppHeader />
      <main class="flex-1 overflow-y-auto p-4 md:p-6">
        <RouterView />
      </main>
    </div>
  </div>
  <RouterView v-else />
</template>
