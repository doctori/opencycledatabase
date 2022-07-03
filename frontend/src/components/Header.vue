<template>
<div>
    <v-app-bar app color="white">
          <v-app-bar-nav-icon></v-app-bar-nav-icon>
          <v-toolbar-title>OCD</v-toolbar-title>

      <v-spacer></v-spacer>

      <v-menu
        left
        bottom
      >
        <template v-slot:activator="{ on, attrs }">
          <v-btn
            icon
            v-bind="attrs"
            v-on="on"
          >
            <v-icon>mdi-dots-vertical</v-icon>
          </v-btn>
        </template>
            <v-list>
          <v-list-item >
            <router-link :to="{name: 'Profile'}">
            <v-list-item-title>Login</v-list-item-title>
            </router-link>
          </v-list-item>
        </v-list>
      </v-menu>
    </v-app-bar>

</div>
</template>

<script>
export default {
  name: 'Header',
  methods: {},
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
