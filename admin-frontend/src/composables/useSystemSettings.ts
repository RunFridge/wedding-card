import { ref } from 'vue';
import api from '@/lib/axios';
import type {
  SystemSettings,
  SystemSettingsUpdateResponse,
  S3TestResult,
} from '@/types/admin';

export function useSystemSettings() {
  const settings = ref<SystemSettings | null>(null);
  const loading = ref(false);
  const saving = ref(false);
  const testing = ref(false);
  const error = ref('');

  async function load() {
    loading.value = true;
    error.value = '';
    try {
      const { data } = await api.get<SystemSettings>('/admin/system-settings');
      settings.value = data;
    } catch {
      error.value = 'Failed to load system settings';
    } finally {
      loading.value = false;
    }
  }

  async function save(
    updates: Partial<SystemSettings>,
  ): Promise<SystemSettingsUpdateResponse | null> {
    saving.value = true;
    error.value = '';
    try {
      const { data } = await api.put<SystemSettingsUpdateResponse>(
        '/admin/system-settings',
        updates,
      );
      await load();
      return data;
    } catch {
      error.value = 'Failed to save system settings';
      return null;
    } finally {
      saving.value = false;
    }
  }

  async function testS3(
    config: Partial<SystemSettings>,
  ): Promise<S3TestResult | null> {
    testing.value = true;
    error.value = '';
    try {
      const { data } = await api.post<S3TestResult>(
        '/admin/system-settings/test-s3',
        config,
      );
      return data;
    } catch {
      error.value = 'Failed to test S3 connection';
      return null;
    } finally {
      testing.value = false;
    }
  }

  async function testModeration(config: {
    openai_api_key?: string;
  }): Promise<S3TestResult | null> {
    testing.value = true;
    error.value = '';
    try {
      const { data } = await api.post<S3TestResult>(
        '/admin/system-settings/test-moderation',
        config,
      );
      return data;
    } catch {
      error.value = 'Failed to test moderation API';
      return null;
    } finally {
      testing.value = false;
    }
  }

  async function restart(): Promise<boolean> {
    try {
      await api.post('/admin/restart');
      return true;
    } catch {
      error.value = 'Failed to restart server';
      return false;
    }
  }

  return {
    settings,
    loading,
    saving,
    testing,
    error,
    load,
    save,
    testS3,
    testModeration,
    restart,
  };
}
