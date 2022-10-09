import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'
import Register from '../views/register/Register.vue'
// import Login from '../views/login/Login.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/about',
    name: 'About',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "about" */ '../views/About.vue')
  },
  {
    path: '/register',
    name: 'register',
    component: Register
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('../views/login/Login.vue') // 这样写成函数的好处是懒加载
  },
  {
    path: '/javatosql',
    name: 'javatosql',
    component: () => import('../views/JavaToSql.vue') // 这样写成函数的好处是懒加载
  },
  {
    path: '/daode',
    name: 'daode',
    component: () => import('../views/Daode.vue') // 这样写成函数的好处是懒加载
  },
  {
    path: '/cal',
    name: 'cal',
    component: () => import('../views/Cal.vue') // 这样写成函数的好处是懒加载
  },
  {
    path: '/split',
    name: 'split',
    component: () => import('../views/Split.vue') // 这样写成函数的好处是懒加载
  },
  {
    path: '/game',
    name: 'game',
    component: () => import('../views/Game.vue') // 这样写成函数的好处是懒加载
  },
  {
    path: '/compareFile',
    name: 'compareFile',
    component: () => import('../views/CompareFile.vue') // 这样写成函数的好处是懒加载
  }
]

const router = new VueRouter({
  // mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
