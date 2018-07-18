<template>
    <div class="container">
        <el-form v-if="!$store.state.authUser" :model="form" label-width="120px">
          <el-form-item label="username">
            <el-input v-model="form.username"></el-input>
          </el-form-item>
          <el-form-item label = "password">
            <el-input v-model="form.password" type="password"></el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="login">Login</el-button>
          </el-form-item>
        </el-form>
        <div v-else>
            Hello {{ $store.state.authUser.id }}!
        </div>

    </div>
</template>

<script>
export default {
    data() {
        return {
            form: {
              username: '',
              password: ''
            }
        }
    },
    methods: {
        async login() {
            try {
              await this.$store.dispatch('login', {
                  username: this.form.username,
                  password: this.form.password
              })
              this.form.username = ''
              this.form.password = ''
      
              let path = this.$route.query.path || "/"
              this.$router.replace(path)
            } catch (e) {
              console.error(e);
            }
        },
    }
}
</script>

<style>
.container {
    padding: 100px;
}
.el-input {
  width: 180px;
}
</style>
