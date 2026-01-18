<template>
    <div class="create-account">
        <h3 class="top-header">add account</h3>

        <div class="currency-selector">
            <span class="currency-text">currency:</span>

            <div class="currency-list">
                <button v-for="currency in currencies" :key="currency.code" class="currency btn"
                    :class="{ selected: currency.code == selectedCurrencyCode }"
                    @click="selectedCurrencyCode = currency.code">
                    <span class="small-text currency-code-text">{{ currency.code }}</span>
                </button>
            </div>
        </div>

        <div class="nav-btns">
            <button class="secondary-btn" @click="$emit('cancel')"> cancel </button>
            <button class="primary-btn" @click="createAccount">create</button>
        </div>
    </div>
</template>

<script>
import { fetchCurrencies } from '@/api/api';

export default {
    name: "CreateAccount",

    data() {
        return {
            currencies: [],
            selectedCurrencyCode: "",
        }
    },

    mounted() {
        this.getCurrenicesList()
    },

    methods: {
        async getCurrenicesList() {
            try {
                const resp = await fetchCurrencies();
                this.currencies = resp;
                if (this.currencies.length > 0) {
                    this.selectedCurrencyCode = this.currencies[0].code;
                }
            } catch (error) {
                console.error("failed to fetch currencies")
            }
        },
        createAccount() {
            this.$emit('created', { currency_code: this.selectedCurrencyCode })
        }
    }
}
</script>

<style scoped>
.create-account {
    display: flex;
    flex: 1;
    flex-direction: column;
    max-width: fit-content;
    gap: 10px;
    align-items: center;
}

.top-header {
    flex: 1;
    align-self: stretch;
}

.currency-selector {
    display: flex;
    flex: 1;
    align-self: stretch;
    justify-content: space-between;
    flex-direction: row;
}

.currency-text {
    color: white;
    justify-items: left;

    margin-right: 10px;
}

.currency-list {
    display: flex;
    flex: 1;
    flex-direction: row;
    gap: 10px;
    align-self: stretch;
    max-width: fit-content;
    justify-items: right;

    justify-self: right;
}

.currency {
    border: 1px solid var(--color-text-muted);
    border-radius: 10px;
    padding: 4px;
}

.currency.selected {
    background-color: var(--color-text-muted);
}

.nav-btns {
    display: flex;
    flex: 1;
    flex-direction: row;
    align-self: stretch;
    justify-content: space-between;
    padding: 5px;
}

.cancel-btn {
    border: 1px solid var(--color-text-muted);
    border-radius: 10px;
    text-align: center;
    align-items: center;
    color: var(--color-text-muted);
    padding: 1px 7px;
}

.create-btn {
    background-color: var(--color-accent);
    border-radius: 10px;
    text-align: center;
    align-items: center;
    color: white;
    padding: 1px 7px;
}
</style>
