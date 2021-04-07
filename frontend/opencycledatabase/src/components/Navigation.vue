<template>
<div>
    <v-navigation-drawer app>
      <v-container>
        <v-row class="text-center">
          <v-col>
            <v-img
              :src="logo"
              />
            <h4 class="display-3 font-weight-bold">
              OCD
            </h4>
          </v-col>
          <v-col >
            <p class="subheading font-weight-regular">
              Welcom to the Open Cycling Database,

            </p>
          </v-col>

          <v-col
            class="mb-5"
            cols="12"
          >
            <h2 class="headline font-weight-bold mb-5">
              Standards
            </h2>

          </v-col>

          <v-col
            class="mb-5"
            cols="12"
          >
            <h2 class="headline font-weight-bold mb-5">
              Components
            </h2>
          </v-col>

          <v-col
            class="mb-5"
            cols="12"
          >
            <h2 class="headline font-weight-bold mb-5">
              Brands
            </h2>

          </v-col>
        </v-row>
      </v-container>
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
