<template>
    <ProfileInfo/>
</template>

<script>
// @ is an alias to /src
import ProfileInfo from '@/components/utils/ProfileInfo.vue'

export default {
  name: 'Login',
  data() {
    return {
      user : {}
    }
  },
  methods: {
    printGerard(){
      console.log("GERARD")
    },
    displayProfile(profile){
      console.log("TIMMY")
      profile.then(response => {
        console.log(response.data.access_token)
      }
      )
      

    }
  },
  components: {
    ProfileInfo
  }
}
</script>
