<template>
    <div class="container">
        <!-- Заголовок сайта -->
        <h1 class="site-title">JXGERcorp Banking</h1>

        <div class ="form-box">
            <h2 class="form-title">Login</h2>
            <form @submit.prevent="handleLogin">
                <!-- Поля ввода -->
                <div id="input-group">
                    <input id="username" v-model="username" placeholder="Username" type="text" class="input-field">
                    <input id="pasword" v-model="password" placeholder="Password" type="password" class="input-field">
                </div>

                <!-- Ошибки -->
                 <div v-if="error" class="error">
                    {{ error }}
                 </div>

                <!-- Кнопки -->
                <div class="button-group">
                    <button @click="goRegister" class="btn btn-secondary">Register</button>
                    <button type="submit" class="btn btn-primary">Sign in</button>

                </div>
            </form>
        </div>
    </div>
</template>

<script>
import { loginUser } from '@/api/api';

export default {
    data() {
        return {
            username: '',
            password: '',
            error: null,
        }
    },

    methods: {
        async handleLogin() {
            // clear errors
            this.error = null;

            // check password
            if (this.password.length === 0) {
                this.error = 'password is too short';
                console.log('stupid2');
                return;
            }


            // call API
            try {
                const userData = {
                    username: this.username,
                    password: this.password,
                };

                const response = await loginUser(userData);

                // proceed successfull registration
                console.log('Login successful:', response);
                // clear form
                this.username = '';
                this.password = '';
                this.$router.push('/main'); // автоматический переход на /main
            } catch(err) {
                console.log('Login failed:', err);
                this.error = err.data?.message || 'Login failed. Please try again';
                console.log('Extracted err message:', this.error)
            }
        },
        goRegister() {
            this.$router.push('/register')
        }
    },
};
</script>

<style scoped>
/* Основной контейнер */
.container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    background-color: #f7f9fc;
    font-family: Arial, sans-serif;
}

/* Рамка формы */
.form-box {
    background-color: #ffffff;
    border: 2px solid #2c3e50; /* Темная рамка */
    border-radius: 10px; /* Скругленные углы */
    padding: 20px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1); /* Легкая тень */
    max-width: 400px;
    width: 100%;
    text-align: center;
}

/* Группа полей ввода */
.input-group {
    display: flex;
    flex-direction: column;
    gap: 15px; /* Отступы между полями */
}

/* Группа кнопок */
.button-group {
    display: flex;
    justify-content: space-between;
    margin-top: 20px;
}

/* Ошибка */
.error {
    color: white;
    background-color: #e74c3c;
    margin: 10px 0;
    font-size: 14px;
    text-align: center;
    padding: 8px 12px;
    border-radius: 6px;
    display: inline-block;
    max-width: 100%;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}
</style>