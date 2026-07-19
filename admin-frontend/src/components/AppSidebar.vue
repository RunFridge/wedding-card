<script setup lang="ts">
import { useRoute } from 'vue-router';
import { watch } from 'vue';
import { useI18n } from 'vue-i18n';
import {
  LayoutDashboard,
  BookOpen,
  Trophy,
  Award,
  Camera,
  Images,
  Settings,
  Wrench,
  KeyRound,
  ScrollText,
  Info,
} from 'lucide-vue-next';
import { Separator } from '@/components/ui/separator';
import { cn } from '@/lib/utils';
import { useSidebar } from '@/composables/useSidebar';

const route = useRoute();
const { t } = useI18n();
const { collapsed, mobileOpen, isMobile, closeMobile } = useSidebar();

watch(
  () => route.path,
  () => {
    if (isMobile.value) closeMobile();
  },
);

const navItems = [
  {
    name: 'dashboard',
    labelKey: 'nav.dashboard',
    icon: LayoutDashboard,
    to: '/',
  },
  {
    name: 'guestbook',
    labelKey: 'nav.guestbook',
    icon: BookOpen,
    to: '/guestbook',
  },
  { name: 'rankings', labelKey: 'nav.rankings', icon: Trophy, to: '/rankings' },
  {
    name: 'hall-of-fame',
    labelKey: 'nav.hallOfFame',
    icon: Award,
    to: '/hall-of-fame',
  },
  { name: 'photos', labelKey: 'nav.photos', icon: Camera, to: '/photos' },
  {
    name: 'asset-photos',
    labelKey: 'nav.assetPhotos',
    icon: Images,
    to: '/asset-photos',
  },
  {
    name: 'settings',
    labelKey: 'nav.settings',
    icon: Settings,
    to: '/settings',
  },
  { name: 'system', labelKey: 'nav.system', icon: Wrench, to: '/system' },
  {
    name: 'change-password',
    labelKey: 'nav.password',
    icon: KeyRound,
    to: '/change-password',
  },
  { name: 'logs', labelKey: 'nav.logs', icon: ScrollText, to: '/logs' },
  { name: 'about', labelKey: 'nav.about', icon: Info, to: '/about' },
];
</script>

<template>
  <!-- Mobile: overlay backdrop -->
  <Teleport to="body">
    <div
      v-if="isMobile && mobileOpen"
      class="fixed inset-0 z-40 bg-black/50"
      @click="closeMobile"
    />
  </Teleport>

  <!-- Sidebar -->
  <aside
    :class="
      cn(
        'flex flex-col border-r bg-background transition-transform duration-200 z-50',
        isMobile
          ? 'fixed inset-y-0 left-0 w-56 shadow-lg ' +
              (mobileOpen ? 'translate-x-0' : '-translate-x-full')
          : 'relative transition-[width] ' + (collapsed ? 'w-14' : 'w-56'),
      )
    "
  >
    <div
      class="flex h-14 items-center overflow-hidden px-4 font-semibold whitespace-nowrap"
    >
      <span v-if="isMobile || !collapsed">{{ t('nav.adminPanel') }}</span>
    </div>
    <Separator />
    <nav class="flex-1 space-y-1 p-2">
      <RouterLink
        v-for="item in navItems"
        :key="item.name"
        :to="item.to"
        :data-tour="`nav-${item.name}`"
        :title="!isMobile && collapsed ? t(item.labelKey) : undefined"
        :class="
          cn(
            'flex items-center rounded-md py-2 text-sm transition-colors overflow-hidden whitespace-nowrap',
            !isMobile && collapsed ? 'justify-center px-0' : 'gap-3 px-3',
            route.name === item.name
              ? 'bg-primary text-primary-foreground'
              : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground',
          )
        "
      >
        <component :is="item.icon" class="h-4 w-4 shrink-0" />
        <span v-if="isMobile || !collapsed">{{ t(item.labelKey) }}</span>
      </RouterLink>
    </nav>
  </aside>
</template>
