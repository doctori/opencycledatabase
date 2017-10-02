import Vue from 'vue'
import Router from 'vue-router'
import Hello from '@/components/Hello'
import Standards from '@/components/Standards'
import Standard from '@/components/Standard'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Hello',
      component: Hello
    },
    {
      path: '/standards',
      name: 'Standards',
      component: Standards,
      children: [
        {
          path: ':id',
          name: 'Standard',
          component: Standard
        }
      ]
    }
  ]
})
