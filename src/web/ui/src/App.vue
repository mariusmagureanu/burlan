<template>

  <div id="app" style="height: 100%">

        <div id="nav">
              <router-link v-if="authenticated" to="/login" v-on:click.native="logout()" replace>Logout</router-link>
        </div>
        <router-view @authenticated="setAuthenticated" />


  </div>
</template>

<script>
import messaging from '@/components/messaging.vue'
import login from '@/components/login.vue'

export default {
  name: "app",
  components: {
    messaging
  },
  data() {
    return {
      authenticated: false,              
    }
  },

  mounted() {
      if(!this.authenticated) {
           this.$router.replace({ name: "login" }).catch(err => {});
      }
  },

  methods: {
      setAuthenticated(status) {
         this.authenticated = status;
      },
      logout() {
         this.authenticated = false;  
         localStorage.removeItem('user-token')      
      },
  },
}
</script>

