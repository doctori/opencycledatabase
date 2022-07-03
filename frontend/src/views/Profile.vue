<template>
    <LoginDialog @logged="this.displayProfile" @gerard="this.printGerard"/>
    <ProfileInfo v-if="user"/>
</template>

<script>
// @ is an alias to /src
import LoginDialog from '@/components/utils/LoginDialog.vue'
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
    LoginDialog,
    ProfileInfo
  }
}
</script>
