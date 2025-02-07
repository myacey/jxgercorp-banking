<template>
    <div class="container">
        <!-- Заголовок сайта -->
        <h1 class="site-title">JXGERcorp Banking</h1>

        <div class ="form-box">
            <h2 class="form-title">Register</h2>
            <form @submit.prevent="handleRegister">
                <!-- Поля ввода -->
                <div id="input-group">
                    <input id="username" v-model="username" placeholder="Username" type="text" class="input-field">
                    <input id="email" v-model="email" placeholder="Email" type="email" class="input-field">
                    <input id="pasword" v-model="password" placeholder="Password" type="password" class="input-field">
                    <input id="configm-pasword" v-model="confirmPassword" placeholder="Confirm Password" type="password" class="input-field">
                </div>

                <!-- Ошибки -->
                 <div v-if="error" class="error">
                    {{ error }}
                 </div>

                <!-- Кнопки -->
                <div class="button-group">
                    <button type="back" class="btn btn-secondary">Back</button>
                    <button type="submit" class="btn btn-primary">Next</button>

                </div>
            </form>
        </div>
    </div>
</template>

<script>
import { registerUser } from '@/api/api';

export default {
    data() {
        return {
            username: '',
            email: '',
            password: '',
            confirmPassword: '',
            error: null,
        }
    },

    methods: {
        async handleRegister() {
            // clear errors
            this.error = null;

            // validate email
            const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            if (!emailRegex.test(this.email)) {
                this.error = 'Invalid email format';
                return;
            }

            // check password
            if (this.password !== this.confirmPassword) {
                this.error = 'Passwords do not match';
                return;
            }

            if (this.password.length < 8) {
                this.error = 'Password too short';
                return;
            }

            // call API
            try {
                const userData = {
                    email: this.email,
                    username: this.username,
                    password: this.password,
                };

                const response = await registerUser(userData);

                // proceed successfull registration
                console.log('Register successful:', response);
                alert('Registration succssful!');

                // clear form
                this.email = '';
                this.username = '';
                this.password = '';
                this.confirmPassword = '';

                this.$router.push('/login'); // автоматический переход на /login
            } catch(err) {
                console.log('Registration failed:', err);
                this.error = err.data?.error?.message || 'Registration failed. Please try again';
                console.log('Extracted err message:', this.error)
            }
        },

        handleBack() {
            console.log('Go back')
        },
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
    color: red;
    margin: 10px 0;
    font-size: 14px;
    text-align: left;
}
</style>