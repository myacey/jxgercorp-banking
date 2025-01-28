import Vue from 'vue'
import VueRouter from 'vue-router'
import HomeView from '../views/HomeView.vue'
import Login from '../views/LoginView.vue'
// import Profile from '../views/ProfileView.vue'
import Register from '../views/RegisterView.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'home',
    component: HomeView
  },
  { path: '/register', component: Register },
  { path: '/login', component: Login }
  // { path: '/profile', component: Profile }
]

const router = new VueRouter({
  routes
})

export default router
