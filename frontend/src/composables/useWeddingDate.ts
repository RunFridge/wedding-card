import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { WEDDING_DATETIME } from '../config/wedding';

export function useWeddingDate() {
  const { locale } = useI18n();

  const date = computed(() => {
    const formatted = WEDDING_DATETIME.toLocaleDateString(locale.value, {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      timeZone: 'Asia/Seoul',
    });
    return locale.value === 'ko' ? formatted.replace(/\. /g, '. ') : formatted;
  });

  const day = computed(() =>
    WEDDING_DATETIME.toLocaleDateString(locale.value, {
      weekday: 'long',
      timeZone: 'Asia/Seoul',
    }),
  );

  const time = computed(() =>
    WEDDING_DATETIME.toLocaleTimeString(locale.value, {
      hour: 'numeric',
      minute: '2-digit',
      hour12: true,
      timeZone: 'Asia/Seoul',
    }),
  );

  return { date, day, time };
}
