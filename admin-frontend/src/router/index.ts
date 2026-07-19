import { createRouter, createWebHistory } from 'vue-router';
import { useAuth } from '@/lib/auth';

const router = createRouter({
  history: createWebHistory('/-/admin/'),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { guest: true },
    },
    {
      path: '/setup',
      name: 'setup',
      component: () => import('@/views/SetupView.vue'),
      meta: { auth: true },
    },
    {
      path: '/',
      name: 'dashboard',
      component: () => import('@/views/DashboardView.vue'),
      meta: { auth: true },
    },
    {
      path: '/guestbook',
      name: 'guestbook',
      component: () => import('@/views/GuestbookView.vue'),
      meta: { auth: true },
    },
    {
      path: '/rankings',
      name: 'rankings',
      component: () => import('@/views/RankingsView.vue'),
      meta: { auth: true },
    },
    {
      path: '/hall-of-fame',
      name: 'hall-of-fame',
      component: () => import('@/views/HallOfFameView.vue'),
      meta: { auth: true },
    },
    {
      path: '/photos',
      name: 'photos',
      component: () => import('@/views/PhotosView.vue'),
      meta: { auth: true },
    },
    {
      path: '/asset-photos',
      name: 'asset-photos',
      component: () => import('@/views/AssetPhotosView.vue'),
      meta: { auth: true },
    },
    {
      path: '/settings',
      name: 'settings',
      component: () => import('@/views/SettingsView.vue'),
      meta: { auth: true },
    },
    {
      path: '/system',
      name: 'system',
      component: () => import('@/views/SystemSettingsView.vue'),
      meta: { auth: true },
    },
    {
      path: '/change-password',
      name: 'change-password',
      component: () => import('@/views/ChangePasswordView.vue'),
      meta: { auth: true },
    },
    {
      path: '/logs',
      name: 'logs',
      component: () => import('@/views/LogsView.vue'),
      meta: { auth: true },
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('@/views/AboutView.vue'),
      meta: { auth: true },
    },
  ],
});

router.beforeEach(async (to) => {
  const { authenticated, setupRequired, checkSession } = useAuth();

  if (to.meta.auth) {
    if (!authenticated.value) {
      const valid = await checkSession();
      if (!valid) return { name: 'login', query: { redirect: to.fullPath } };
    }

    if (setupRequired.value && to.name !== 'setup') {
      return { name: 'setup' };
    }

    if (!setupRequired.value && to.name === 'setup') {
      return { name: 'dashboard' };
    }
  }

  if (to.meta.guest && authenticated.value) {
    if (setupRequired.value) return { name: 'setup' };
    return { name: 'dashboard' };
  }
});

export default router;
