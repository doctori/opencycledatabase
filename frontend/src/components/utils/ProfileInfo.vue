<template>
    <div>
        <h2>You are : {{ user.login }} </h2>
    <div>
        Your token is {{ user.token }} <br/>
        Country : {{ user.country }}
    </div>
    </div>
</template>

<script>
//import BackendApiClient from '../../services/BackendApiClient'
export default {
    name: 'ProfileInfo',
    props: {
        user
    },
    components: {},
    data(){
        return {},
    }
}
</script>

