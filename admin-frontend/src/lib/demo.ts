import { ref } from 'vue';

export const isDemo = ref(false);

let loaded = false;

export async function loadDemoFlag(): Promise<void> {
  if (loaded) return;
  loaded = true;
  try {
    const res = await fetch('/api/health');
    const data = await res.json();
    isDemo.value = data.demo === true;
  } catch {
    isDemo.value = false;
  }
}
