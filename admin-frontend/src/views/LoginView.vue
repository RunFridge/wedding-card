<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { useAuth } from '@/lib/auth';
import api from '@/lib/axios';
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Info } from 'lucide-vue-next';

const route = useRoute();
const router = useRouter();
const { t } = useI18n();
const { login } = useAuth();

const password = ref('');
const error = ref('');
const loading = ref(false);
const setupRequired = ref(false);

onMounted(async () => {
  try {
    const res = await api.get('/health');
    setupRequired.value = !!res.data.setup_required;
  } catch {
    // ignore
  }
});

async function handleLogin() {
  if (!password.value) return;
  error.value = '';
  loading.value = true;
  try {
    const ok = await login(password.value);
    if (ok) {
      const redirect =
        typeof route.query.redirect === 'string' ? route.query.redirect : null;
      router.push(redirect || { name: 'dashboard' });
    } else {
      error.value = t('login.invalidPassword');
    }
  } catch {
    error.value = t('login.invalidPassword');
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div class="flex min-h-screen items-center justify-center">
    <Card class="w-full max-w-sm">
      <CardHeader>
        <CardTitle class="text-center text-2xl">{{
          t('login.title')
        }}</CardTitle>
        <CardDescription v-if="setupRequired" class="text-center">
          {{ t('login.setupRequired') }}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div
          v-if="setupRequired"
          class="flex items-start gap-2 rounded-md border bg-muted px-4 py-3 text-sm text-muted-foreground mb-4"
        >
          <Info class="h-4 w-4 mt-0.5 shrink-0" />
          <div>
            <p>{{ t('login.initialPasswordGenerated') }}</p>
            <i18n-t keypath="login.findPassword" tag="p" class="mt-1">
              <template #file>
                <code
                  class="rounded bg-background px-1 py-0.5 text-xs font-mono"
                  >admin_password.txt</code
                >
              </template>
            </i18n-t>
            <p class="mt-1 text-xs">
              Docker:
              <code class="rounded bg-background px-1 py-0.5 font-mono"
                >docker compose exec wedding-server cat
                /data/admin_password.txt</code
              >
            </p>
          </div>
        </div>
        <form class="space-y-4" @submit.prevent="handleLogin">
          <Input
            v-model="password"
            type="password"
            :placeholder="t('login.passwordPlaceholder')"
            autofocus
          />
          <p v-if="error" class="text-sm text-destructive">{{ error }}</p>
          <Button type="submit" class="w-full" :disabled="loading">
            {{ loading ? t('login.loggingIn') : t('login.login') }}
          </Button>
        </form>
      </CardContent>
    </Card>
  </div>
</template>
