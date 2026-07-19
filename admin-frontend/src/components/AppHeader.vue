<script setup lang="ts">
import { onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import {
  LogOut,
  PanelLeftClose,
  PanelLeftOpen,
  Menu,
  Users,
  ExternalLink,
  Languages,
  CircleHelp,
} from 'lucide-vue-next';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { useAuth } from '@/lib/auth';
import { useSidebar } from '@/composables/useSidebar';
import { useAdminWS } from '@/composables/useAdminWS';
import { startTour, startTourOnFirstVisit } from '@/composables/useAdminTour';
import { messages, setLocale, type Locale } from '@/i18n';

const router = useRouter();
const { t } = useI18n();
const { logout } = useAuth();
const { collapsed, isMobile, toggle } = useSidebar();
const { visitorCount, connected } = useAdminWS();
const locales = Object.keys(messages) as Locale[];

async function handleLogout() {
  await logout();
  router.push({ name: 'login' });
}

onMounted(startTourOnFirstVisit);
</script>

<template>
  <header class="flex h-14 items-center justify-between border-b px-4 md:px-6">
    <Button variant="ghost" size="icon" @click="toggle">
      <Menu v-if="isMobile" class="h-4 w-4" />
      <PanelLeftOpen v-else-if="collapsed" class="h-4 w-4" />
      <PanelLeftClose v-else class="h-4 w-4" />
    </Button>
    <div class="flex items-center gap-2">
      <div
        data-tour="visitors"
        class="flex items-center gap-1.5 text-sm text-muted-foreground"
        :title="
          connected ? t('header.wsConnected') : t('header.wsDisconnected')
        "
      >
        <span class="relative flex h-2 w-2">
          <span
            v-if="connected"
            class="absolute inline-flex h-full w-full animate-ping rounded-full bg-green-400 opacity-75"
          ></span>
          <span
            class="relative inline-flex h-2 w-2 rounded-full"
            :class="connected ? 'bg-green-500' : 'bg-muted-foreground/40'"
          ></span>
        </span>
        <Users class="h-4 w-4" />
        <span>{{ visitorCount }}</span>
      </div>
      <Button variant="ghost" size="sm" as="a" href="/" target="_blank">
        <ExternalLink class="mr-2 h-4 w-4" />
        {{ t('header.home') }}
      </Button>
      <Button
        variant="ghost"
        size="icon"
        data-tour="help"
        :title="t('header.showGuide')"
        @click="startTour"
      >
        <CircleHelp class="h-4 w-4" />
      </Button>
      <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <Button variant="ghost" size="icon" :title="t('header.language')">
            <Languages class="h-4 w-4" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuItem
            v-for="locale in locales"
            :key="locale"
            @click="setLocale(locale)"
          >
            {{ messages[locale].app.localeLabel }}
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
      <Button variant="ghost" size="sm" @click="handleLogout">
        <LogOut class="mr-2 h-4 w-4" />
        {{ t('header.logout') }}
      </Button>
    </div>
  </header>
</template>
