<script setup lang="ts">
import { onMounted, reactive, ref, watch, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { useConfig } from '@/composables/useConfig';
import type { WeddingConfig, BusInfoEntry, MapProviders } from '@/types/admin';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Save, Plus, X, ExternalLink } from 'lucide-vue-next';

function defaultMapProviders(): MapProviders {
  return {
    embed_provider: '',
    latitude: 0,
    longitude: 0,
    api_key: '',
    links: { google: '', kakao: '', naver: '', tmap: '' },
  };
}

const EMBED_PROVIDER_OPTIONS = [
  { value: 'none', labelKey: 'settings.providerNone' },
  { value: 'google', labelKey: 'settings.googleMaps' },
  { value: 'kakao', labelKey: 'settings.kakaoMap' },
  { value: 'naver', labelKey: 'settings.naverMap' },
  { value: 'tmap', labelKey: 'settings.tmap' },
] as const;

const MAP_PROVIDER_DOCS: Record<string, string> = {
  google: 'https://developers.google.com/maps/documentation/embed/get-started',
  kakao: 'https://apis.map.kakao.com/web/guide/',
  naver: 'https://navermaps.github.io/maps.js.ncp/docs/',
  tmap: 'https://tmapapi.tmapmobility.com/main.html',
};

const { t } = useI18n();
const { config, loading, saving, error, load, save } = useConfig();

const VALID_TABS = [
  'general',
  'bank',
  'venue',
  'map',
  'transport',
  'game-moderation',
] as const;
const route = useRoute();
const router = useRouter();
const activeTab = computed({
  get: () => {
    const hash = (route.hash || '').replace('#', '');
    return VALID_TABS.includes(hash as any) ? hash : 'general';
  },
  set: (v: string) => {
    router.replace({ hash: `#${v}` });
  },
});

const form = reactive<WeddingConfig>({
  groom_eng_name: '',
  groom_kor_name: '',
  bride_eng_name: '',
  bride_kor_name: '',
  groom_father_kor_name: '',
  groom_mother_kor_name: '',
  bride_father_kor_name: '',
  bride_mother_kor_name: '',
  groom_bank_account: '',
  bride_bank_account: '',
  groom_father_bank_account: '',
  groom_mother_bank_account: '',
  bride_father_bank_account: '',
  bride_mother_bank_account: '',
  wedding_datetime: '',
  venue_name: '',
  venue_address: '',
  venue_floor: '',
  venue_hall: '',
  venue_phone: '',
  map_providers: defaultMapProviders(),
  subway_info: [],
  bus_info: [],
  charter_bus: [] as {
    location: string;
    company: string;
    bus_number: string;
    departure: string;
  }[],
  charter_bus_notice: '',
  car_info: '',
  groom_birth_order: '',
  bride_birth_order: '',
  card_game_timer: 30000,
  game_npc_message: '',
  avatar_colors: '',
  short_greeting: '',
  main_greet_text: '',
  simple_redirect_url: '',
  photo_upload_enabled: false,
  photo_upload_hours_before: 1,
  hearts_flush_interval_ms: 2000,
  hearts_flush_batch_size: 50,
  moderation_thresholds: {} as Record<string, number>,
});

const YEARS = Array.from({ length: 7 }, (_, i) => 2024 + i);
const HOURS = Array.from({ length: 24 }, (_, i) => i);
const MINUTES = Array.from({ length: 12 }, (_, i) => i * 5);
const TZ_OPTIONS = [
  { value: '-12:00', label: 'UTC-12:00' },
  { value: '-11:00', label: 'UTC-11:00' },
  { value: '-10:00', label: 'UTC-10:00 (Hawaii)' },
  { value: '-09:00', label: 'UTC-09:00 (Alaska)' },
  { value: '-08:00', label: 'UTC-08:00 (PST)' },
  { value: '-07:00', label: 'UTC-07:00 (MST)' },
  { value: '-06:00', label: 'UTC-06:00 (CST)' },
  { value: '-05:00', label: 'UTC-05:00 (EST)' },
  { value: '-04:00', label: 'UTC-04:00' },
  { value: '-03:00', label: 'UTC-03:00' },
  { value: '-02:00', label: 'UTC-02:00' },
  { value: '-01:00', label: 'UTC-01:00' },
  { value: '+00:00', label: 'UTC+00:00 (GMT)' },
  { value: '+01:00', label: 'UTC+01:00 (CET)' },
  { value: '+02:00', label: 'UTC+02:00 (EET)' },
  { value: '+03:00', label: 'UTC+03:00' },
  { value: '+04:00', label: 'UTC+04:00' },
  { value: '+05:00', label: 'UTC+05:00' },
  { value: '+05:30', label: 'UTC+05:30 (IST)' },
  { value: '+06:00', label: 'UTC+06:00' },
  { value: '+07:00', label: 'UTC+07:00' },
  { value: '+08:00', label: 'UTC+08:00' },
  { value: '+09:00', label: 'UTC+09:00 (KST/JST)' },
  { value: '+10:00', label: 'UTC+10:00 (AEST)' },
  { value: '+11:00', label: 'UTC+11:00' },
  { value: '+12:00', label: 'UTC+12:00 (NZST)' },
];

const dtYear = ref(2026);
const dtMonth = ref(6);
const dtDay = ref(6);
const dtHour = ref(11);
const dtMinute = ref(0);
const dtTz = ref('+09:00');

const daysInMonth = computed(() =>
  new Date(dtYear.value, dtMonth.value, 0).getDate(),
);

function parseDatetime(iso: string) {
  const m = iso.match(
    /^(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2}):\d{2}([+-]\d{2}:\d{2})$/,
  );
  if (m) {
    dtYear.value = parseInt(m[1]);
    dtMonth.value = parseInt(m[2]);
    dtDay.value = parseInt(m[3]);
    dtHour.value = parseInt(m[4]);
    const rawMin = Math.round(parseInt(m[5]) / 5) * 5;
    dtMinute.value = rawMin >= 60 ? 55 : rawMin;
    dtTz.value = m[6];
  }
}

watch([dtYear, dtMonth, dtDay, dtHour, dtMinute, dtTz], () => {
  if (dtDay.value > daysInMonth.value) dtDay.value = daysInMonth.value;
  const mo = String(dtMonth.value).padStart(2, '0');
  const d = String(dtDay.value).padStart(2, '0');
  const h = String(dtHour.value).padStart(2, '0');
  const mi = String(dtMinute.value).padStart(2, '0');
  form.wedding_datetime = `${dtYear.value}-${mo}-${d}T${h}:${mi}:00${dtTz.value}`;
});

const success = ref(false);

function populateForm(cfg: WeddingConfig) {
  const mp = cfg.map_providers ?? defaultMapProviders();
  Object.assign(form, {
    ...cfg,
    map_providers: {
      ...defaultMapProviders(),
      ...mp,
      links: { ...defaultMapProviders().links, ...mp.links },
    },
    subway_info: [...cfg.subway_info],
    bus_info: cfg.bus_info.map((b: BusInfoEntry) => ({ ...b })),
    charter_bus: (cfg.charter_bus || []).map((b: any) => ({ ...b })),
    moderation_thresholds: { ...cfg.moderation_thresholds },
  });
  parseDatetime(form.wedding_datetime);
}

onMounted(async () => {
  await load();
  if (config.value) populateForm(config.value);
});

function addSubwayLine() {
  form.subway_info.push('');
}

function removeSubwayLine(i: number) {
  form.subway_info.splice(i, 1);
}

function addBusEntry() {
  form.bus_info.push({ stop: '', routes: '' });
}

function removeBusEntry(i: number) {
  form.bus_info.splice(i, 1);
}

function addCharterBus() {
  form.charter_bus.push({
    location: '',
    company: '',
    bus_number: '',
    departure: '',
  });
}

function removeCharterBus(i: number) {
  form.charter_bus.splice(i, 1);
}

async function handleSave() {
  success.value = false;
  const toSave: WeddingConfig = {
    ...form,
  };
  try {
    await save(toSave);
    success.value = true;
    setTimeout(() => (success.value = false), 3000);
  } catch {
    // error is set by composable
  }
}
</script>

<template>
  <div>
    <div class="mb-4 flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold">{{ t('nav.settings') }}</h1>
        <p class="text-sm text-muted-foreground">
          {{ t('settings.subtitle') }}
        </p>
      </div>
    </div>

    <div v-if="loading" class="space-y-4">
      <Skeleton v-for="i in 4" :key="i" class="h-40 w-full" />
    </div>

    <div v-else>
      <Tabs
        :model-value="activeTab"
        @update:model-value="(v) => (activeTab = String(v))"
      >
        <TabsList class="mb-6 flex-wrap">
          <TabsTrigger value="general">{{
            t('settings.tabGeneral')
          }}</TabsTrigger>
          <TabsTrigger value="bank">{{ t('settings.tabBank') }}</TabsTrigger>
          <TabsTrigger value="venue">{{ t('settings.tabVenue') }}</TabsTrigger>
          <TabsTrigger value="map">{{ t('settings.tabMap') }}</TabsTrigger>
          <TabsTrigger value="transport">{{
            t('settings.tabTransport')
          }}</TabsTrigger>
          <TabsTrigger value="game-moderation">{{
            t('settings.tabGame')
          }}</TabsTrigger>
        </TabsList>

        <!-- General: Names + Wedding Details -->
        <TabsContent value="general" class="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>{{ t('settings.names') }}</CardTitle>
            </CardHeader>
            <CardContent>
              <div class="grid gap-4 sm:grid-cols-2">
                <div class="space-y-4">
                  <p class="text-sm font-semibold text-muted-foreground">
                    {{ t('settings.groomSide') }}
                  </p>
                  <div>
                    <label class="mb-1 block text-sm font-medium">{{
                      t('settings.englishName')
                    }}</label>
                    <Input v-model="form.groom_eng_name" placeholder="Groom" />
                  </div>
                  <div>
                    <label class="mb-1 block text-sm font-medium">{{
                      t('settings.koreanName')
                    }}</label>
                    <Input v-model="form.groom_kor_name" placeholder="김철수" />
                  </div>
                  <div>
                    <label class="mb-1 block text-sm font-medium">{{
                      t('settings.fatherName')
                    }}</label>
                    <Input
                      v-model="form.groom_father_kor_name"
                      placeholder="김아버지"
                    />
                  </div>
                  <div>
                    <label class="mb-1 block text-sm font-medium">{{
                      t('settings.motherName')
                    }}</label>
                    <Input
                      v-model="form.groom_mother_kor_name"
                      placeholder="박어머니"
                    />
                  </div>
                  <div>
                    <label class="mb-1 block text-sm font-medium">{{
                      t('settings.birthOrder')
                    }}</label>
                    <Input
                      v-model="form.groom_birth_order"
                      placeholder="장남"
                    />
                  </div>
                </div>
                <div class="space-y-4">
                  <p class="text-sm font-semibold text-muted-foreground">
                    {{ t('settings.brideSide') }}
                  </p>
                  <div>
                    <label class="mb-1 block text-sm font-medium">{{
                      t('settings.englishName')
                    }}</label>
                    <Input v-model="form.bride_eng_name" placeholder="Bride" />
                  </div>
                  <div>
                    <label class="mb-1 block text-sm font-medium">{{
                      t('settings.koreanName')
                    }}</label>
                    <Input v-model="form.bride_kor_name" placeholder="이영희" />
                  </div>
                  <div>
                    <label class="mb-1 block text-sm font-medium">{{
                      t('settings.fatherName')
                    }}</label>
                    <Input
                      v-model="form.bride_father_kor_name"
                      placeholder="이아버지"
                    />
                  </div>
                  <div>
                    <label class="mb-1 block text-sm font-medium">{{
                      t('settings.motherName')
                    }}</label>
                    <Input
                      v-model="form.bride_mother_kor_name"
                      placeholder="최어머니"
                    />
                  </div>
                  <div>
                    <label class="mb-1 block text-sm font-medium">{{
                      t('settings.birthOrder')
                    }}</label>
                    <Input
                      v-model="form.bride_birth_order"
                      placeholder="장녀"
                    />
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>{{ t('settings.weddingDetails') }}</CardTitle>
            </CardHeader>
            <CardContent class="grid gap-4">
              <div>
                <label class="mb-2 block text-sm font-medium">{{
                  t('settings.dateTime')
                }}</label>
                <div class="space-y-2">
                  <div class="flex flex-wrap items-center gap-2">
                    <Select
                      :model-value="String(dtYear)"
                      @update:model-value="(v) => (dtYear = Number(v))"
                    >
                      <SelectTrigger class="w-24"
                        ><SelectValue
                      /></SelectTrigger>
                      <SelectContent>
                        <SelectItem
                          v-for="y in YEARS"
                          :key="y"
                          :value="String(y)"
                          >{{ y }}</SelectItem
                        >
                      </SelectContent>
                    </Select>
                    <span class="text-muted-foreground">/</span>
                    <Select
                      :model-value="String(dtMonth)"
                      @update:model-value="(v) => (dtMonth = Number(v))"
                    >
                      <SelectTrigger class="w-20"
                        ><SelectValue
                      /></SelectTrigger>
                      <SelectContent>
                        <SelectItem
                          v-for="m in 12"
                          :key="m"
                          :value="String(m)"
                          >{{ String(m).padStart(2, '0') }}</SelectItem
                        >
                      </SelectContent>
                    </Select>
                    <span class="text-muted-foreground">/</span>
                    <Select
                      :model-value="String(dtDay)"
                      @update:model-value="(v) => (dtDay = Number(v))"
                    >
                      <SelectTrigger class="w-20"
                        ><SelectValue
                      /></SelectTrigger>
                      <SelectContent>
                        <SelectItem
                          v-for="d in daysInMonth"
                          :key="d"
                          :value="String(d)"
                          >{{ String(d).padStart(2, '0') }}</SelectItem
                        >
                      </SelectContent>
                    </Select>
                  </div>
                  <div class="flex flex-wrap items-center gap-2">
                    <Select
                      :model-value="String(dtHour)"
                      @update:model-value="(v) => (dtHour = Number(v))"
                    >
                      <SelectTrigger class="w-20"
                        ><SelectValue
                      /></SelectTrigger>
                      <SelectContent>
                        <SelectItem
                          v-for="h in HOURS"
                          :key="h"
                          :value="String(h)"
                          >{{ String(h).padStart(2, '0') }}</SelectItem
                        >
                      </SelectContent>
                    </Select>
                    <span class="text-muted-foreground">:</span>
                    <Select
                      :model-value="String(dtMinute)"
                      @update:model-value="(v) => (dtMinute = Number(v))"
                    >
                      <SelectTrigger class="w-20"
                        ><SelectValue
                      /></SelectTrigger>
                      <SelectContent>
                        <SelectItem
                          v-for="m in MINUTES"
                          :key="m"
                          :value="String(m)"
                          >{{ String(m).padStart(2, '0') }}</SelectItem
                        >
                      </SelectContent>
                    </Select>
                    <Select
                      :model-value="dtTz"
                      @update:model-value="(v) => (dtTz = String(v))"
                    >
                      <SelectTrigger class="w-52"
                        ><SelectValue
                      /></SelectTrigger>
                      <SelectContent>
                        <SelectItem
                          v-for="tz in TZ_OPTIONS"
                          :key="tz.value"
                          :value="tz.value"
                          >{{ tz.label }}</SelectItem
                        >
                      </SelectContent>
                    </Select>
                  </div>
                </div>
              </div>
              <div>
                <label class="mb-1 block text-sm font-medium">{{
                  t('settings.shortGreeting')
                }}</label>
                <Input
                  v-model="form.short_greeting"
                  placeholder="저희 결혼합니다"
                />
                <p class="mt-1 text-xs text-muted-foreground">
                  {{ t('settings.shortGreetingHint') }}
                </p>
              </div>
              <div>
                <label class="mb-1 block text-sm font-medium">{{
                  t('settings.greetingText')
                }}</label>
                <textarea
                  v-model="form.main_greet_text"
                  rows="3"
                  placeholder="소중한 분들을 초대합니다.&#10;함께 축복해 주시면 더없는 기쁨으로 간직하겠습니다."
                  class="flex w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                />
              </div>
              <div>
                <label class="mb-1 block text-sm font-medium">{{
                  t('settings.simpleRedirectUrl')
                }}</label>
                <Input
                  v-model="form.simple_redirect_url"
                  placeholder="https://..."
                />
                <p class="mt-1 text-xs text-muted-foreground">
                  {{ t('settings.simpleRedirectUrlHint') }}
                </p>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <!-- Bank Accounts -->
        <TabsContent value="bank" class="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>{{ t('settings.bankAccounts') }}</CardTitle>
            </CardHeader>
            <CardContent>
              <div class="grid gap-4 sm:grid-cols-2">
                <div class="space-y-4">
                  <p class="text-sm font-semibold text-muted-foreground">
                    {{ t('settings.groomSide') }}
                  </p>
                  <div>
                    <label class="mb-1 block text-sm font-medium">{{
                      t('settings.groom')
                    }}</label>
                    <Input
                      v-model="form.groom_bank_account"
                      placeholder="카카오뱅크 0000-00-0000000"
                    />
                  </div>
                  <div>
                    <label class="mb-1 block text-sm font-medium">{{
                      t('settings.groomFather')
                    }}</label>
                    <Input
                      v-model="form.groom_father_bank_account"
                      placeholder="은행 000-000000-00000"
                    />
                  </div>
                  <div>
                    <label class="mb-1 block text-sm font-medium">{{
                      t('settings.groomMother')
                    }}</label>
                    <Input
                      v-model="form.groom_mother_bank_account"
                      placeholder="은행 000-000000-00000"
                    />
                  </div>
                </div>
                <div class="space-y-4">
                  <p class="text-sm font-semibold text-muted-foreground">
                    {{ t('settings.brideSide') }}
                  </p>
                  <div>
                    <label class="mb-1 block text-sm font-medium">{{
                      t('settings.bride')
                    }}</label>
                    <Input
                      v-model="form.bride_bank_account"
                      placeholder="은행 000-000000-00000"
                    />
                  </div>
                  <div>
                    <label class="mb-1 block text-sm font-medium">{{
                      t('settings.brideFather')
                    }}</label>
                    <Input
                      v-model="form.bride_father_bank_account"
                      placeholder="은행 000-000000-00000"
                    />
                  </div>
                  <div>
                    <label class="mb-1 block text-sm font-medium">{{
                      t('settings.brideMother')
                    }}</label>
                    <Input
                      v-model="form.bride_mother_bank_account"
                      placeholder="은행 000-000000-00000"
                    />
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <!-- Venue -->
        <TabsContent value="venue" class="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>{{ t('settings.venue') }}</CardTitle>
            </CardHeader>
            <CardContent class="grid gap-4 sm:grid-cols-2">
              <div class="sm:col-span-2">
                <label class="mb-1 block text-sm font-medium">{{
                  t('settings.venueName')
                }}</label>
                <Input
                  v-model="form.venue_name"
                  placeholder="OO웨딩홀"
                />
              </div>
              <div class="sm:col-span-2">
                <label class="mb-1 block text-sm font-medium">{{
                  t('settings.address')
                }}</label>
                <Input
                  v-model="form.venue_address"
                  placeholder="서울특별시 중구 세종대로 110"
                />
              </div>
              <div>
                <label class="mb-1 block text-sm font-medium">{{
                  t('settings.floor')
                }}</label>
                <Input v-model="form.venue_floor" placeholder="15층" />
              </div>
              <div>
                <label class="mb-1 block text-sm font-medium">{{
                  t('settings.hall')
                }}</label>
                <Input v-model="form.venue_hall" placeholder="그랜드홀" />
              </div>
              <div>
                <label class="mb-1 block text-sm font-medium">{{
                  t('settings.phone')
                }}</label>
                <Input v-model="form.venue_phone" placeholder="02-000-0000" />
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <!-- Map -->
        <TabsContent value="map" class="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>{{ t('settings.embed') }}</CardTitle>
              <p class="text-sm text-muted-foreground">
                {{ t('settings.embedHint') }}
                <template
                  v-if="
                    form.map_providers.embed_provider &&
                    MAP_PROVIDER_DOCS[form.map_providers.embed_provider]
                  "
                >
                  <a
                    :href="MAP_PROVIDER_DOCS[form.map_providers.embed_provider]"
                    target="_blank"
                    class="inline-flex items-center gap-1 text-foreground hover:underline"
                  >
                    {{ t('settings.apiDocs') }} <ExternalLink class="h-3 w-3" />
                  </a>
                </template>
              </p>
            </CardHeader>
            <CardContent class="space-y-4">
              <div>
                <label class="mb-1 block text-sm font-medium">{{
                  t('settings.provider')
                }}</label>
                <Select
                  :model-value="form.map_providers.embed_provider || 'none'"
                  @update:model-value="
                    (v) =>
                      (form.map_providers.embed_provider =
                        String(v) === 'none' ? '' : String(v))
                  "
                >
                  <SelectTrigger class="w-48">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem
                      v-for="opt in EMBED_PROVIDER_OPTIONS"
                      :key="opt.value"
                      :value="opt.value"
                    >
                      {{ t(opt.labelKey) }}
                    </SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <template v-if="form.map_providers.embed_provider">
                <div class="grid gap-4 sm:grid-cols-2">
                  <div>
                    <label class="mb-1 block text-sm font-medium">{{
                      t('settings.latitude')
                    }}</label>
                    <Input
                      v-model.number="form.map_providers.latitude"
                      type="number"
                      step="0.0001"
                      placeholder="37.5665"
                    />
                  </div>
                  <div>
                    <label class="mb-1 block text-sm font-medium">{{
                      t('settings.longitude')
                    }}</label>
                    <Input
                      v-model.number="form.map_providers.longitude"
                      type="number"
                      step="0.0001"
                      placeholder="126.9780"
                    />
                  </div>
                </div>
                <div>
                  <label class="mb-1 block text-sm font-medium">{{
                    t('settings.apiKey')
                  }}</label>
                  <Input
                    v-model="form.map_providers.api_key"
                    placeholder="API key for embed"
                  />
                </div>
              </template>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>{{ t('settings.links') }}</CardTitle>
              <p class="text-sm text-muted-foreground">
                {{ t('settings.linksHint') }}
              </p>
            </CardHeader>
            <CardContent class="grid gap-4">
              <div>
                <label class="mb-1 block text-sm font-medium">{{
                  t('settings.googleMaps')
                }}</label>
                <Input
                  v-model="form.map_providers.links.google"
                  placeholder="https://maps.app.goo.gl/..."
                />
              </div>
              <div>
                <label class="mb-1 block text-sm font-medium">{{
                  t('settings.kakaoMap')
                }}</label>
                <Input
                  v-model="form.map_providers.links.kakao"
                  placeholder="https://kko.to/..."
                />
              </div>
              <div>
                <label class="mb-1 block text-sm font-medium">{{
                  t('settings.naverMap')
                }}</label>
                <Input
                  v-model="form.map_providers.links.naver"
                  placeholder="https://naver.me/..."
                />
              </div>
              <div>
                <label class="mb-1 block text-sm font-medium">{{
                  t('settings.tmap')
                }}</label>
                <Input
                  v-model="form.map_providers.links.tmap"
                  placeholder="https://tmap.life/..."
                />
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <!-- Transport -->
        <TabsContent value="transport" class="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>{{ t('settings.transportInfo') }}</CardTitle>
            </CardHeader>
            <CardContent class="space-y-6">
              <div>
                <div class="mb-2 flex items-center justify-between">
                  <label class="text-sm font-medium">{{
                    t('settings.subwayInfo')
                  }}</label>
                  <Button variant="outline" size="sm" @click="addSubwayLine">
                    <Plus class="mr-1 h-3 w-3" /> {{ t('settings.add') }}
                  </Button>
                </div>
                <div
                  v-for="(_, i) in form.subway_info"
                  :key="i"
                  class="mb-2 flex gap-2"
                >
                  <Input
                    v-model="form.subway_info[i]"
                    class="flex-1"
                    placeholder="1호선 시청역 1번 출구"
                  />
                  <Button
                    variant="ghost"
                    size="icon"
                    class="h-9 w-9 shrink-0"
                    @click="removeSubwayLine(i)"
                  >
                    <X class="h-4 w-4" />
                  </Button>
                </div>
              </div>

              <Separator />

              <div>
                <div class="mb-2 flex items-center justify-between">
                  <label class="text-sm font-medium">{{
                    t('settings.busInfo')
                  }}</label>
                  <Button variant="outline" size="sm" @click="addBusEntry">
                    <Plus class="mr-1 h-3 w-3" /> {{ t('settings.add') }}
                  </Button>
                </div>
                <div
                  v-for="(bus, i) in form.bus_info"
                  :key="i"
                  class="mb-3 rounded-md border p-3"
                >
                  <div class="mb-2 flex items-start justify-between">
                    <Badge variant="outline">{{
                      t('settings.stopN', { n: i + 1 })
                    }}</Badge>
                    <Button
                      variant="ghost"
                      size="icon"
                      class="h-7 w-7"
                      @click="removeBusEntry(i)"
                    >
                      <X class="h-4 w-4" />
                    </Button>
                  </div>
                  <div class="space-y-2">
                    <Input v-model="bus.stop" placeholder="Bus stop name" />
                    <Input
                      v-model="bus.routes"
                      placeholder="Routes (e.g. 13-1, 20-2, 202)"
                    />
                  </div>
                </div>
              </div>

              <Separator />

              <div>
                <div class="mb-2 flex items-center justify-between">
                  <label class="text-sm font-medium">{{
                    t('settings.charterBus')
                  }}</label>
                  <Button variant="outline" size="sm" @click="addCharterBus">
                    <Plus class="mr-1 h-3 w-3" /> {{ t('settings.add') }}
                  </Button>
                </div>
                <div
                  v-for="(cb, i) in form.charter_bus"
                  :key="i"
                  class="mb-3 rounded-md border p-3"
                >
                  <div class="mb-2 flex items-start justify-between">
                    <Badge variant="outline">{{
                      t('settings.busN', { n: i + 1 })
                    }}</Badge>
                    <Button
                      variant="ghost"
                      size="icon"
                      class="h-7 w-7"
                      @click="removeCharterBus(i)"
                    >
                      <X class="h-4 w-4" />
                    </Button>
                  </div>
                  <div class="space-y-2">
                    <Input
                      v-model="cb.location"
                      placeholder="Location (e.g. OO경기장 주차장)"
                    />
                    <Input
                      v-model="cb.company"
                      placeholder="Company (e.g. OO관광)"
                    />
                    <Input
                      v-model="cb.bus_number"
                      placeholder="Bus number (e.g. 00가0000)"
                    />
                    <Input
                      v-model="cb.departure"
                      placeholder="Departure time (e.g. 7시 30분 출발)"
                    />
                  </div>
                </div>
              </div>

              <div>
                <label class="mb-1 block text-sm font-medium">{{
                  t('settings.charterBusNotice')
                }}</label>
                <textarea
                  v-model="form.charter_bus_notice"
                  rows="4"
                  placeholder="전세버스 탑승 안내 문구를 입력하세요"
                  class="flex w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                />
                <p class="mt-1 text-xs text-muted-foreground">
                  {{ t('settings.charterBusNoticeHint') }}
                </p>
              </div>

              <Separator />

              <div>
                <label class="mb-1 block text-sm font-medium">{{
                  t('settings.carParkingInfo')
                }}</label>
                <Input
                  v-model="form.car_info"
                  placeholder="전용 주차장 이용 안내"
                />
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <!-- Game: Game & Display + Hearts Settings -->
        <TabsContent value="game-moderation" class="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>{{ t('settings.gameDisplay') }}</CardTitle>
            </CardHeader>
            <CardContent class="grid gap-4 sm:grid-cols-2">
              <div>
                <label class="mb-1 block text-sm font-medium">{{
                  t('settings.cardGameTimer')
                }}</label>
                <Input
                  v-model.number="form.card_game_timer"
                  type="number"
                  min="1000"
                  max="120000"
                  step="1000"
                  placeholder="30000"
                />
                <p class="mt-1 text-xs text-muted-foreground">
                  {{ t('settings.cardGameTimerHint') }}
                </p>
              </div>
              <div class="sm:col-span-2">
                <label class="mb-1 block text-sm font-medium">{{
                  t('settings.gameNpcMessage')
                }}</label>
                <Input
                  v-model="form.game_npc_message"
                  placeholder="NPC intro message for first-time players"
                />
                <p class="mt-1 text-xs text-muted-foreground">
                  {{ t('settings.gameNpcMessageHint') }}
                </p>
              </div>
              <div>
                <label class="mb-1 block text-sm font-medium">{{
                  t('settings.avatarColors')
                }}</label>
                <Input
                  v-model="form.avatar_colors"
                  placeholder="8B6914,A0722A,C4943A,5C3A0E,F5E6C8"
                />
                <p class="mt-1 text-xs text-muted-foreground">
                  {{ t('settings.avatarColorsHint') }}
                </p>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>{{ t('settings.photoUpload') }}</CardTitle>
            </CardHeader>
            <CardContent class="grid gap-4 sm:grid-cols-2">
              <div>
                <label class="mb-1 block text-sm font-medium">{{
                  t('settings.hoursBeforeWedding')
                }}</label>
                <Input
                  v-model.number="form.photo_upload_hours_before"
                  type="number"
                  min="0"
                  max="720"
                  step="0.5"
                  placeholder="1"
                />
                <p class="mt-1 text-xs text-muted-foreground">
                  {{ t('settings.hoursBeforeWeddingHint') }}
                </p>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>{{ t('settings.heartsSettings') }}</CardTitle>
            </CardHeader>
            <CardContent class="grid gap-4 sm:grid-cols-2">
              <div>
                <label class="mb-1 block text-sm font-medium">{{
                  t('settings.flushInterval')
                }}</label>
                <Input
                  v-model.number="form.hearts_flush_interval_ms"
                  type="number"
                  min="500"
                  max="30000"
                  step="500"
                />
                <p class="mt-1 text-xs text-muted-foreground">
                  {{ t('settings.flushIntervalHint') }}
                </p>
              </div>
              <div>
                <label class="mb-1 block text-sm font-medium">{{
                  t('settings.flushBatchSize')
                }}</label>
                <Input
                  v-model.number="form.hearts_flush_batch_size"
                  type="number"
                  min="1"
                  max="1000"
                  step="10"
                />
                <p class="mt-1 text-xs text-muted-foreground">
                  {{ t('settings.flushBatchSizeHint') }}
                </p>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      <div class="mt-6 flex items-center gap-3 pb-8">
        <Button :disabled="saving" @click="handleSave">
          <Save class="mr-2 h-4 w-4" />
          {{ saving ? t('common.saving') : t('common.saveSettings') }}
        </Button>
        <span v-if="success" class="text-sm text-green-600">{{
          t('common.settingsSaved')
        }}</span>
        <span v-if="error" class="text-sm text-destructive">{{ error }}</span>
      </div>
    </div>
  </div>
</template>
