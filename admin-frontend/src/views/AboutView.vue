<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { Github, ExternalLink } from 'lucide-vue-next';

const { t } = useI18n();
const version = ref('');
const loading = ref(true);

onMounted(async () => {
  try {
    const res = await fetch('/api/health');
    const data = await res.json();
    version.value = data.version ?? 'unknown';
  } catch {
    version.value = 'unknown';
  } finally {
    loading.value = false;
  }
});

interface Dependency {
  name: string;
  url: string;
  description: string;
}

const backendDeps: Dependency[] = [
  {
    name: 'chi',
    url: 'https://github.com/go-chi/chi',
    description: 'HTTP router',
  },
  {
    name: 'httprate',
    url: 'https://github.com/go-chi/httprate',
    description: 'Rate limiting',
  },
  {
    name: 'gorilla/websocket',
    url: 'https://github.com/gorilla/websocket',
    description: 'WebSocket',
  },
  {
    name: 'go-sqlite3',
    url: 'https://github.com/mattn/go-sqlite3',
    description: 'SQLite driver',
  },
  {
    name: 'aws-sdk-go-v2',
    url: 'https://github.com/aws/aws-sdk-go-v2',
    description: 'AWS S3 SDK',
  },
  {
    name: 'golang.org/x/crypto',
    url: 'https://pkg.go.dev/golang.org/x/crypto',
    description: 'Cryptography',
  },
  {
    name: 'golang.org/x/image',
    url: 'https://pkg.go.dev/golang.org/x/image',
    description: 'Image processing',
  },
];

const visitorDeps: Dependency[] = [
  { name: 'Vue.js', url: 'https://vuejs.org', description: 'UI framework' },
  {
    name: 'Vue Router',
    url: 'https://router.vuejs.org',
    description: 'Routing',
  },
  {
    name: 'Cropper.js',
    url: 'https://fengyuanchen.github.io/cropperjs/',
    description: 'Image cropping',
  },
  { name: 'Phaser', url: 'https://phaser.io', description: 'Game engine' },
  {
    name: 'Tailwind CSS',
    url: 'https://tailwindcss.com',
    description: 'CSS framework',
  },
  { name: 'Vite', url: 'https://vite.dev', description: 'Build tool' },
];

const adminDeps: Dependency[] = [
  { name: 'Vue.js', url: 'https://vuejs.org', description: 'UI framework' },
  {
    name: 'Vue Router',
    url: 'https://router.vuejs.org',
    description: 'Routing',
  },
  {
    name: 'Reka UI',
    url: 'https://reka-ui.com',
    description: 'Headless UI primitives',
  },
  { name: 'Lucide', url: 'https://lucide.dev', description: 'Icons' },
  { name: 'Axios', url: 'https://axios-http.com', description: 'HTTP client' },
  {
    name: 'class-variance-authority',
    url: 'https://cva.style',
    description: 'Variant styling',
  },
  {
    name: 'Tailwind Merge',
    url: 'https://github.com/dcastil/tailwind-merge',
    description: 'Tailwind class merging',
  },
  { name: 'VueUse', url: 'https://vueuse.org', description: 'Vue composables' },
  {
    name: 'Tailwind CSS',
    url: 'https://tailwindcss.com',
    description: 'CSS framework',
  },
  { name: 'Vite', url: 'https://vite.dev', description: 'Build tool' },
];

const sections = [
  { titleKey: 'about.backend', deps: backendDeps },
  { titleKey: 'about.visitorFrontend', deps: visitorDeps },
  { titleKey: 'about.adminFrontend', deps: adminDeps },
];
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold">{{ t('nav.about') }}</h1>
    <p class="mb-6 text-sm text-muted-foreground">{{ t('about.subtitle') }}</p>

    <div class="space-y-4">
      <Card>
        <CardHeader>
          <CardTitle class="text-sm font-medium text-muted-foreground">{{
            t('about.version')
          }}</CardTitle>
        </CardHeader>
        <CardContent>
          <Skeleton v-if="loading" class="h-8 w-24" />
          <div v-else class="text-3xl font-bold">v{{ version }}</div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle class="text-sm font-medium text-muted-foreground">{{
            t('about.sourceCode')
          }}</CardTitle>
        </CardHeader>
        <CardContent>
          <a
            href="https://github.com/RunFridge/wedding-card"
            target="_blank"
            rel="noopener noreferrer"
            class="inline-flex items-center gap-2 text-sm font-medium text-primary hover:underline"
          >
            <Github class="h-4 w-4" />
            RunFridge/wedding-card
            <ExternalLink class="h-3 w-3" />
          </a>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle class="text-sm font-medium text-muted-foreground">{{
            t('about.attribution')
          }}</CardTitle>
        </CardHeader>
        <CardContent class="space-y-6">
          <div v-for="section in sections" :key="section.titleKey">
            <h3 class="mb-2 text-sm font-semibold">
              {{ t(section.titleKey) }}
            </h3>
            <ul class="space-y-1">
              <li v-for="dep in section.deps" :key="dep.name" class="text-sm">
                <a
                  :href="dep.url"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="font-medium text-primary hover:underline"
                  >{{ dep.name }}</a
                >
                <span class="text-muted-foreground">
                  — {{ dep.description }}</span
                >
              </li>
            </ul>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
