<template>
    <div class="account-selector">
        <div v-for="account in accounts" :key="account.id" class="account-item"
            :class="{ active: account.id == selectedAccountID }" @click="$emit('select', account)">
            <div class="date-and-id">
                <p class="small-text created-at">{{ formatDate(account.created_at) }}</p>
                <p class="small-text account-id">{{ account.id }}</p>
            </div>

            <p class="divider">|</p>

            <div class="balance-with-currency">
                <h3 class="balance">{{ account.balance }}</h3>
                <div class="small-text currency">{{ account.currency }}</div>
            </div>
        </div>

        <!-- КНОПКА ДОБАВЛЕНИЯ НОВОГО АККАУНТА -->
        <div class="account-item add-account">
            +
        </div>
    </div>
</template>

<script>
export default {
    name: "AccountSelector",

    props: {
        accounts: Array,
        selectedAccountID: String
    },

    mounted() {
        console.log(this.accounts)
    },

    methods: {
        formatDate(date) {
            return new Date(date).toLocaleDateString('ru-RU')
        }
    }
}
</script>

<style scoped>
.account-selector {
    position: absolute;
    top: 0;
    right: 100%;

    margin-right: 12px;
    z-index: 1001;

    display: inline-flex;
    max-width: fit-content;
    flex-direction: column;
    align-items: flex-end;
    padding: 4px;
    gap: 10px;
    border-radius: 10px;

    background-color: var(--color-bg-alt)
}

.account-item {
    display: flex;
    flex: 1;
    flex-direction: row;
    justify-content: space-between;
    gap: 2px;
    padding: 5px;
    align-self: stretch;

    cursor: pointer;
}

.account-item:hover {
    opacity: 0.8;
}

.account-item.active {
    border: 1px solid var(--color-text-muted);
    border-radius: 10px;
}

.date-and-id {
    display: flex;
    flex: 1;
    flex-direction: column;
    justify-content: left;
    align-self: flex-start;
    margin-right: 10px;
}

.created-at {
    text-align: right;
    color: var(--color-text-muted)
}

.account-id {
    text-align: right;
    color: var(--color-text-muted);
    white-space: nowrap;
}

.divider {
    align-self: center;
    white-space: nowrap;
}

.balance-with-currency {
    display: flex;
    align-items: center;
}

.add-account {
    display: flex;
    flex: 1;
    align-self: stretch;
    border: 1px solid var(--color-text-muted);
    border-radius: 10px;

    justify-content: center;
    color: var(--color-text-muted)
}
</style>
