<template>
<div>
    <v-navigation-drawer>
      <v-list nav>
          <v-list-item>
            <v-img
              :src="logo"
              />
            <router-link
            :to="{ name: 'Home'}"
            >
              <h4 class="display-3 font-weight-bold">
                OCD
              </h4>
            </router-link>
          </v-list-item>
          <v-list-item >
            <p class="subheading font-weight-regular">
              {{ $t("messages.welcome_nav") }}

            </p>
          </v-list-item>

          <v-list-item
            class="mb-5"
            cols="12"
          >
          <router-link :to="{name:'Standards'}">
            <h2 class="headline font-weight-bold mb-5">
              Standards
            </h2>
          </router-link>

          </v-list-item>

          <v-list-item
            class="mb-5"
            cols="12"
          >
          <router-link
            :to="{name: 'Components'}">
            <h2 class="headline font-weight-bold mb-5">
              {{ $t('messages.components')}}
            </h2>
          </router-link>
          </v-list-item>

          <v-list-item
            class="mb-5"
            cols="12"
          >
          <router-link
            :to="{name:'Brands'}"
            >
            <h2 class="headline font-weight-bold mb-5">
              Brands
            </h2>
          </router-link>
          </v-list-item>
      </v-list>
    </v-navigation-drawer>
</div>
</template>

<script>
import logo from '../assets/logo.png';
export default {
  name: 'Navigation',

  data() {    
    return{
      logo,
    }}
}
</script>
