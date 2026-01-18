<template>
  <div class="container">
    <p>confirmation status: {{ emailConfirmedStatus }}</p>
  </div>
</template>

<script>
import { confirmUserEmail } from '@/api/api'
import { defineComponent } from 'vue'

export default defineComponent({
  data() {
    return {
      emailConfirmedStatus: 'invalid'
    }
  },

  mounted() {
    console.log('AAAAA;;;')
    this.confirmUsrEmail()
  },

  methods: {
    async confirmUsrEmail() {
      const { username, code } = this.$route.query
      if (!username || !code) {
        console.error('Missing query parameters')
        return
      }
      console.log('username:', username)
      console.log('code:', code)

      const params = {
        username: username,
        code: code
      }

      const resp = await confirmUserEmail(params)
      this.emailConfirmedStatus = resp.message
    }
  }
})
</script>

<style scoped>
.container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background-color: #f7f9fc;
  font-family: Arial, sans-serif;
}
</style>
