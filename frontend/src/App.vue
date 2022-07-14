<template>
  <v-app id="ocd">
    <Header/>
    <Navigation/>
    <v-main>
    <router-view/>
    </v-main>
    <Footer/>
  </v-app>
</template>

<script>
import Navigation from './components/Navigation';
import Header from './components/Header';
import Footer from './components/Footer';
export default {
  name: 'App',
  components: {
    Navigation,
    Header,
    Footer
  },
  data() {
    return {
    }
  },
  provide() {
    return {}
  }
}
</script>
