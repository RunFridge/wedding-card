import { describe, it, expect } from 'vitest';
import { ref } from 'vue';
import { usePagination } from '../composables/usePagination';

function makeItems(n: number) {
  return ref(Array.from({ length: n }, (_, i) => i));
}

describe('usePagination', () => {
  it('returns first page of items', () => {
    const items = makeItems(25);
    const { paginatedItems } = usePagination(items, 10);
    expect(paginatedItems.value).toHaveLength(10);
    expect(paginatedItems.value[0]).toBe(0);
  });

  it('computes totalPages correctly', () => {
    const items = makeItems(25);
    const { totalPages } = usePagination(items, 10);
    expect(totalPages.value).toBe(3);
  });

  it('returns minimum 1 total page for empty list', () => {
    const items = makeItems(0);
    const { totalPages } = usePagination(items, 10);
    expect(totalPages.value).toBe(1);
  });

  it('navigates to page 2', () => {
    const items = makeItems(25);
    const { page, paginatedItems } = usePagination(items, 10);
    page.value = 2;
    expect(paginatedItems.value).toHaveLength(10);
    expect(paginatedItems.value[0]).toBe(10);
  });

  it('last page has remaining items', () => {
    const items = makeItems(25);
    const { page, paginatedItems } = usePagination(items, 10);
    page.value = 3;
    expect(paginatedItems.value).toHaveLength(5);
  });

  it('setPageSize resets to page 1', () => {
    const items = makeItems(50);
    const { page, pageSize, setPageSize } = usePagination(items, 10);
    page.value = 3;
    setPageSize(25);
    expect(page.value).toBe(1);
    expect(pageSize.value).toBe(25);
  });

  it('resetPage goes back to page 1', () => {
    const items = makeItems(30);
    const { page, resetPage } = usePagination(items, 10);
    page.value = 3;
    resetPage();
    expect(page.value).toBe(1);
  });

  it('reacts to items changing', () => {
    const items = ref([1, 2, 3]);
    const { totalPages, paginatedItems } = usePagination(items, 2);
    expect(totalPages.value).toBe(2);
    items.value = [1, 2, 3, 4, 5];
    expect(totalPages.value).toBe(3);
    expect(paginatedItems.value).toHaveLength(2);
  });
});
