<template>

  <div id="app" style="height: 100%">

        <div id="nav">
              <router-link v-if="authenticated" to="/login" v-on:click.native="logout()" replace>Logout</router-link>
        </div>
        <router-view @authenticated="setAuthenticated" />


  </div>
</template>

<script>
import UsersList from '@/components/UsersList.vue'
import login from '@/components/login.vue'

export default {
  name: "app",
  components: {
    UsersList
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
    async getEmployees() {
      const myHeaders = new Headers();
      myHeaders.append('Content-Type', 'application/json');
      try {
        const response = await fetch('http://localhost:8080/api/v1/users', {
         method: 'GET',
        })
        const data = await response.json()
        console.log(data)
        this.employees = data
      } catch (error) {
        console.error(error)
      }
    },

    async addEmployee(employee) {
      try {
        const response = await fetch('https://jsonplaceholder.typicode.com/users', {
          method: 'POST',
          body: JSON.stringify(employee),
          headers: { "Content-type": "application/json; charset=UTF-8" }
        })
        const data = await response.json()
        this.employees = [...this.employees, data]
      } catch (error) {
        console.error(error)
      }
    },

    async editEmployee(id, updatedEmployee) {
      try {
        const response = await fetch(`https://jsonplaceholder.typicode.com/users/${id}`, {
          method: 'PUT',
          body: JSON.stringify(updatedEmployee),
          headers: { "Content-type": "application/json; charset=UTF-8" }
        })
        const data = await response.json()
        this.employees = this.employees.map(employee => employee.id === id ? data : employee)
      } catch (error) {
        console.error(error)
      }
    },

    async deleteEmployee(id) {
      try {
        await fetch(`https://jsonplaceholder.typicode.com/users/${id}`, {
          method: 'DELETE'
        })
        this.employees = this.employees.filter(employee => employee.id !== id)
      } catch (error) {
        console.error(error)
      }
    },
  },
}
</script>

