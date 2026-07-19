import { createRouter, createWebHistory } from 'vue-router';
import HomeView from '../views/HomeView.vue';

const routes = [
  { path: '/', name: 'home', component: HomeView },
  { path: '/guestbook', name: 'guestbook', component: () => import('../views/GuestbookView.vue') },
  { path: '/game', name: 'game', component: () => import('../views/GameView.vue') },
  { path: '/map', name: 'map', component: () => import('../views/MapView.vue') },
  { path: '/photo', name: 'photo', component: () => import('../views/PhotoView.vue') },
  { path: '/achievements', name: 'achievements', component: () => import('../views/AchievementsView.vue') },
  {
    path: '/hall-of-fame',
    name: 'hall-of-fame',
    component: () => import('../views/HallOfFameView.vue'),
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => import('../views/NotFoundView.vue'),
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 };
  },
});

export default router;
