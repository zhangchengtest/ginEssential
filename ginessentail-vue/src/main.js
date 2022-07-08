import Vuelidate from 'vuelidate'
import { BootstrapVue, IconsPlugin } from 'bootstrap-vue'
import axios from 'axios'
import VueAxios from 'vue-axios'
import Vue from 'vue'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import locale from 'element-ui/lib/locale/lang/zh-CN'

import App from './App.vue'
import router from './router'
import store from './store'

// scss style
import './assets/scss/index.scss'

Vue.config.productionTip = false
Vue.config.devtools = true

Vue.use(BootstrapVue)
Vue.use(IconsPlugin)
Vue.use(Vuelidate)
Vue.use(VueAxios, axios)

Vue.use(ElementUI, { locale })

new Vue({
  router,
  store,
  render: (h) => h(App)
}).$mount('#app')
