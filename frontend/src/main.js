import { createApp } from 'vue'
import { useI18n } from 'vue-i18n'
import App from './App.vue'

// vuetify
import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import AuthService from '@/services/AuthService'
import router from '@/router'
import i18n from './i18n'

const vuetify = createVuetify({
  components,
  directives
})

const app = createApp(App)
  .use(AuthService)
  .use(vuetify)
  .use(router)
  .use(i18n)
  .mount('#app')
