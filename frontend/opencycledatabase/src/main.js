import 'mdb-vue-ui-kit/css/mdb.min.css'
import { createApp } from 'vue'
import axios from 'axios'
import App from './App.vue'

var app = createApp(App)
app.mount('#app')

app.config.globalProperties.axios=axios
