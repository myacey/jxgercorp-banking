import { getUsernameFromToken } from '@/utils/auth';
import Vue from 'vue';
import VueRouter from 'vue-router';

Vue.use(VueRouter)

const routes = [
  { path: '/login', component: () => import('@/views/LoginView.vue') },
  { path: '/register', component: () => import('@/views/RegisterView.vue') },
  { path: '/main', component: () => import('@/views/MainView.vue') },
  { path: '/user/confirm', component: () => import('@/views/ConfirmView.vue') },
  { path: '/', redirect: '/main' } // При заходе на '/' идём на '/main'
]

const router = new VueRouter({
  routes
});

// Глобальная проверка токена перед каждой навигацией
router.beforeEach(async (to, from, next) => {
  const username = getUsernameFromToken(); // Проверяем наличие username в токене

  if (!username && to.path !== '/login' && to.path !== '/register' && to.path !== '/user/confirm') {
      next('/login'); // Если username нет → на login
  } else if (username && (to.path === '/login' || to.path === '/register')) {
      next('/main'); // Если username есть, но юзер на login → перекидываем на main
  } else {
      next(); // Иначе даём перейти
  }
});


export default router;
