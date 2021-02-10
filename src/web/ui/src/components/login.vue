<template>
    <div id="login">
        <h1>Login</h1>
        <input type="text" name="username" v-model="input.username" placeholder="Username" />
        <input type="password" name="password" v-model="input.password" placeholder="Password" />
        <button type="button" v-on:click="login()">Login</button>
    </div>
</template>

<script>
    export default {
        name: 'Login',
        data() {
            return {
                input: {
                    username: "",
                    password: ""
                }
            }
        },
        methods: {
            async login() {
                if(this.input.username != "" && this.input.password != "") {
                    const response = await fetch('http://localhost:8080/api/v1/login/'+this.input.username, {
                               method: 'POST',
                              })

                    const data = await response.headers.get('X-Jwt')
                    console.log(data)

                    if(response.status == 200) {
                        localStorage.setItem('user-token', data)
                        this.$emit("authenticated", true);
                        this.$router.replace({ name: "users" });
                    } else {
                        console.log("The username and / or password is incorrect");
                    }
                } else {
                    console.log("A username and password must be present");
                }
            }
        }
    }
</script>

<style scoped>
    #login {
        width: 500px;
        border: 1px solid #CCCCCC;
        background-color: #FFFFFF;
        margin: auto;
        margin-top: 200px;
        padding: 20px;
    }
</style>
