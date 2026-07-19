import { ref, onMounted, onUnmounted } from 'vue';

const collapsed = ref(false);
const mobileOpen = ref(false);
const isMobile = ref(false);

const MOBILE_BREAKPOINT = 768;

function checkMobile() {
  isMobile.value = window.innerWidth < MOBILE_BREAKPOINT;
  if (!isMobile.value) mobileOpen.value = false;
}

export function useSidebar() {
  onMounted(() => {
    checkMobile();
    window.addEventListener('resize', checkMobile);
  });

  onUnmounted(() => {
    window.removeEventListener('resize', checkMobile);
  });

  function toggle() {
    if (isMobile.value) {
      mobileOpen.value = !mobileOpen.value;
    } else {
      collapsed.value = !collapsed.value;
    }
  }

  function closeMobile() {
    mobileOpen.value = false;
  }

  return { collapsed, mobileOpen, isMobile, toggle, closeMobile };
}
