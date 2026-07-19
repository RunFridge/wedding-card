import { ref } from 'vue';

export interface BulkProgress {
  total: number;
  completed: number;
  failed: number;
}

export function useBulkAction() {
  const running = ref(false);
  const progress = ref<BulkProgress | null>(null);

  async function execute(
    ids: number[],
    action: (id: number) => Promise<void>,
    onDone: () => Promise<void>,
  ): Promise<number> {
    running.value = true;
    progress.value = { total: ids.length, completed: 0, failed: 0 };
    let failed = 0;

    for (const id of ids) {
      try {
        await action(id);
      } catch {
        failed++;
      }
      progress.value = {
        total: ids.length,
        completed: progress.value.completed + 1,
        failed,
      };
    }

    await onDone();
    running.value = false;
    progress.value = null;
    return failed;
  }

  return { running, progress, execute };
}
