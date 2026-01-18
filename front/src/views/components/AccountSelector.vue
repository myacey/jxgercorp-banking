<template>
    <div class="account-selector">
        <div v-for="account in accounts" :key="account.id" class="account-item btn"
            :class="{ active: account.id == selectedAccountID }" @click="$emit('select', account)">
            <div class="account-details">

                <div class="date-and-delete">
                    <div class="trashcan" :style="{ maskImage: `url(${trashcanSVG})` }"
                        @click.stop="$emit('request-delete', account)"></div>
                    <p class="small-text created-at">{{ formatDate(account.created_at) }}</p>
                </div>

                <p class="small-text account-id">{{ account.id }}</p>
            </div>

            <p class="divider">|</p>

            <div class="balance-with-currency">
                <h3 class="balance">{{ account.balance }}</h3>
                <div class="small-text currency">{{ account.currency }}</div>
            </div>
        </div>

        <!-- КНОПКА ДОБАВЛЕНИЯ НОВОГО АККАУНТА -->
        <div class="account-item add-account btn" v-if="isCreatingAccount == false" @click="isCreatingAccount = true">
            <span v-if="isCreatingAccount == false" class="small-text"> + </span>

        </div>

        <div class="add-account" v-if="isCreatingAccount == true">
            <CreateAccount @cancel="isCreatingAccount = false" @created="onAccountCreated" />
        </div>

    </div>
</template>

<script>
import CreateAccount from './CreateAccount.vue';

export default {
    name: "AccountSelector",
    components: {
        CreateAccount
    },

    props: {
        accounts: {
            type: Array,
            required: true,
        },
        selectedAccountID: {
            type: String,
            required: false,
        }
    },

    data() {
        return {
            isCreatingAccount: false,
            trashcanSVG: require('@/assets/trashcan.svg')
        }
    },

    mounted() {
        console.log(this.accounts)
    },

    methods: {
        formatDate(date) {
            return new Date(date).toLocaleDateString('ru-RU')
        },
        onAccountCreated(payload) {
            this.$emit('account-created', payload)
            this.isCreatingAccount = false
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

.account-item.active {
    border: 1px solid var(--color-text-muted);
    border-radius: 10px;
}

.account-details {
    display: flex;
    flex: 1;
    flex-direction: column;
    justify-content: left;
    align-self: flex-start;
    margin-right: 10px;
}

.date-and-delete {
    display: flex;
    flex: 1;
    flex-direction: row;
    justify-content: space-between;
}

.trashcan {
    width: 18px;
    height: 18px;
    opacity: 0;
    pointer-events: none;
    transition: opacity 0.4s ease;
    background-color: var(--color-accent);
    mask-size: contain;
    mask-repeat: no-repeat;
    mask-position: center;
    -webkit-mask-size: contain;
    -webkit-mask-repeat: no-repeat;
    -webkit-mask-position: center;
}

.account-item:hover .trashcan {
    opacity: 1;
    pointer-events: auto;
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
