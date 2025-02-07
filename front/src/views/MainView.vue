<template>
    <div class="container">
        <!-- Верхнее меню -->
        <div class="dashboard">
            <div class="header">
                <span class="username">{{ username }}</span>
                <div class="avatar"></div>
            </div>
          </div>
        <!-- Заголовок сайта -->
        <h1 class="site-title">JXGERcorp Banking</h1>

        <!-- Основной контет -->
        <div class="content">
            <!-- Левое меню баланса -->
            <div class="main-box">
                <div class="balance-section">
                    <button class="btn btn-primary" @click="showModal = true">New Transaction</button>
                    <p class="balance">${{ balance }}</p>
                    <span class="owner">{{ username }}</span>
                </div>
            </div>

            <!-- Правое меню истории -->
            <div class="history">
                <h2>history</h2>
                <div class="divider"></div>

                <div v-for="(entry, index) in entries" :key="index" class="entries">
                    <div class="avatar-and-entry-name">
                        <div class="avatar"></div>
                        <span class="entry-name">{{ entry.withUser }}</span>
                    </div>
                    <span :class="entry.amount > 0 ? 'positive' : 'negative'" >
                        {{ entry.amount > 0 ? '+' : '' }}{{ entry.amount }}
                    </span>
                </div>

                <button class="btn btn-secondary">see more -></button>
            </div>
        </div>

        <!-- Окно создания перевода -->
        <div v-if="showModal" class="modal-overlay" @click.self="showModal=false">
            <div class="modal">
                <h3>New Transaction</h3>
                <input v-model="recipient" type="text" placeholder="Recipient">
                <input v-model.number="amount" type="number" placeholder="Amount">

                <div class="buttons">
                    <button class="btn btn-secondary" @click="showModal=false">Cancel</button>
                    <button class="btn btn-primary" @click="commitTransaction">Commit</button>
                </div>
            </div>
        </div>
        
    </div>
</template>

<script>
import { createTransaction, getUserBalance, searchEntries } from '@/api/api';
import { getUsernameFromToken } from '@/utils/auth';
import { convertKeysToCamel } from '@/utils/snake2camel';
import { snake } from "case";

export default {
    data() {
        return {
            balance: 0,
            username: '',
            entries: [],

            showModal: false,
            recipient: '',
            amount: '',
        }
    },
    methods: {
        async getLastEntries() {
            const searchTrxData = {
                offset: 0,
                limit: 5,
            }
            const snakeSearchTrxData = Object.fromEntries(
                Object.entries(searchTrxData).map(([key, value]) => [snake(key), value])
            );
            const responseSnake = await searchEntries(searchTrxData);
            console.log('fetched transactions:', responseSnake);

            this.entries = (convertKeysToCamel(responseSnake));
            console.log('transactions:', this.entries);
        },
        async commitTransaction() {
            if (!this.recipient || !this.amount) return;
            try {
                const trxData = {
                    toUser: this.recipient,
                    amount: this.amount,
                }
                const snakeTrxData = Object.fromEntries(
                    Object.entries(trxData).map(([key, value]) => [snake(key), value])
                );
                await createTransaction(snakeTrxData);
                this.showModal = false;
                this.recipient = '';
                this.amount = '';
                this.getBalance();
                this.getLastEntries();
            } catch(err) {
                console.log('Transaction failed:', err)
            }
        },
        seeMore() {
            this.$router.push('/history');
        },
        async getBalance() {
            try {
                this.balance = await getUserBalance();
            } catch (error) {
                console.log('Cant recieve usr balance:', error)
            }
        }
    },

    mounted() {
        this.getLastEntries();

        this.username = getUsernameFromToken();
        this.getBalance();
        console.log("user info: username:",this.username, "; balance:", this.balance);
    },
};
</script>

<style scoped>
.modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    animation: fadeIn 0.3s;
}

.modal {
    background: white;
    padding: 20px;
    border-radius: 10px;
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
    display: flex;
    flex-direction: column;
    gap: 10px;
    animation: slideIn 0.3s;
}

/* Основной контейнер */
.container {
    display: flex;
    flex-direction: column;
    align-items: center;
    min-height: 100vh;
    background-color: #f7f9fc;
    font-family: Arial, Helvetica, sans-serif;

}

/* Верхняя панель */
.dashboard {
    width: 100%; /* Занимает всю ширину */
    background: white;
    display: flex;
    justify-content: right;
    border-radius: 20px;
    padding: 20px;
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
    align-items: center;
    margin-bottom: 20px;
}

/* Данные в dashboard */
.header {
    display: flex;
    justify-content: space-around;
    gap: 10px;
    align-items: center;
    padding-top: 5px;
    padding-bottom: 5px;
    border-bottom: 2px solid #ddd;
}

/* Никнейм в dashboard */
.username {
    color: #7b3f98;
    font-size: 18px;
    font-weight: bold;
}


/* Аватар */
.avatar {
    width: 30px;
    height: 30px;
    background: black;
    border-radius: 50%;
}

/* Основной контейнер */
.content {
    display: flex;
    width: 50%;
    min-width: 600px;
    max-width: 1200px;
    justify-content: space-evenly;
    background-color: #dedada;
    border-radius: 30px;
    font-family: Arial, Helvetica, sans-serif;
    padding: 30px;
}

/* Левый раздел баланса */
.balance-section {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 20px 0;
    margin: 50px;
}

.balance {
    font-size: 40px;
    font-weight: bold;
    color:#7b3f98;
    margin-top: 10px;
    margin-bottom: 10px;
}

.owner {
    font-size: 14px;
}

/* Правый раздел истории */
.history {
    margin: 50px;
}

.divider {
    height: 2px;
    background: black;
    margin-bottom: 10px;
}

.avatar-and-entry-name {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 2px;
}

.entries {
    display: flex;
    align-items: center;
    gap: 20px;
    justify-content: space-between;
    padding: 10px 0;
    border-bottom: 1px solid #ddd;
}

.positive {
    color: green;
}

.negative {
    color: darkred;
}

</style>
