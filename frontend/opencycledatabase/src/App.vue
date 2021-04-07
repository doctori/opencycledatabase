<template>
  <v-app id="ocd">

    <Header/>
    <Navigation/>

    <!-- Sizes your content based upon application components -->
    <Main/>
    <v-footer>
      <Footer/>
    </v-footer>
  </v-app>
</template>

<script>
import Main from './components/Main.vue';
import Navigation from './components/Navigation.vue';
import Header from './components/Header.vue';
import Footer from './components/Footer.vue';

export default {
  name: 'App',

  components: {
    Main,
    Navigation,
    Header,
    Footer,
  },

  data: () => ({
    //
  }),
};
</script>
