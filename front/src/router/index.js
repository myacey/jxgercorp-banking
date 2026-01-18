import { getUsernameFromToken } from "@/utils/auth";
import Vue from "vue";
import VueRouter from "vue-router";

Vue.use(VueRouter);

const routes = [
  {
    path: "/login",
    component: () => import("@/views/pages/Login.vue"),
  },
  {
    path: "/register",
    component: () => import("@/views/pages/Register.vue"),
  },
  {
    path: "/main",
    component: () => import("@/views/pages/Homepage.vue"),
  },
  {
    path: "/user/confirm",
    component: () => import("@/views/pages/Confirm.vue"),
  },
  {
    path: "/",
    redirect: "/main",
  }, // При заходе на '/' идём на '/main'
];

// const router = new VueRouter({
//   routes
// });

const router = new VueRouter({
  mode: "history",
  routes,
});

// Глобальная проверка токена перед каждой навигацией
router.beforeEach((to, from, next) => {
  const username = getUsernameFromToken();

  const publicPages = ["/login", "/register", "/user/confirm"];
  const authRequired = !publicPages.includes(to.path);

  // если пользователь не авторизован и идёт на приватную страницу
  if (!username && authRequired) {
    if (to.path !== "/login") {
      return next("/login");
    }
  }

  // если пользователь уже авторизован и идёт на публичную страницу
  if (username && publicPages.includes(to.path)) {
    if (to.path !== "/main") {
      return next({ path: "/main", replace: true });
    }
  }

  // если всё ок — просто продолжаем
  next();
});

export default router;
