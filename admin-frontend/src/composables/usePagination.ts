import { ref, computed, type Ref } from 'vue';

export function usePagination<T>(items: Ref<T[]>, initialPageSize = 10) {
  const page = ref(1);
  const pageSize = ref(initialPageSize);

  const totalPages = computed(() =>
    Math.max(1, Math.ceil(items.value.length / pageSize.value)),
  );

  const paginatedItems = computed(() => {
    const start = (page.value - 1) * pageSize.value;
    return items.value.slice(start, start + pageSize.value);
  });

  function setPageSize(size: number) {
    pageSize.value = size;
    page.value = 1;
  }

  function resetPage() {
    page.value = 1;
  }

  return { page, pageSize, totalPages, paginatedItems, setPageSize, resetPage };
}
