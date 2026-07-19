import { i18n } from '@/i18n';

export function apiErrorMessage(e: unknown): string {
  const { t, te } = i18n.global;
  const data = (e as { response?: { data?: unknown } })?.response?.data;
  if (data && typeof data === 'object') {
    const { code, error } = data as { code?: string; error?: string };
    if (code && te(`errors.${code}`)) return t(`errors.${code}`);
    if (error) return error;
  }
  if (typeof data === 'string' && data.trim()) return data.trim();
  return t('errors.requestFailed');
}
