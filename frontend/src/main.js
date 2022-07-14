import Vue from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify';
import router from './router'
import i18n from './i18n'
import AuthService from '@/services/AuthService'
Vue.config.productionTip = false
Vue.use(AuthService)
new Vue({
  vuetify,
  router,
  i18n,
  render: h => h(App)
}).$mount('#app')
