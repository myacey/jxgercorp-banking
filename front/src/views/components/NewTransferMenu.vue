<template>
    <div class="new-transfer">
        <h3 class="header">new transfer</h3>

        <!-- Указание имени пользователя, кому делаем перевод -->
        <div class="username-input-block">
            <input class="username-input base-text" type="text" v-model="username" placeholder="username" />
            <div class="search-btn btn" :class="{ active: newSearchAvailable }"
                @click.stop="$emit('request-search-accounts', username); lastSearchedUsername = username">
                <div class="search-icon" :style="{ maskImage: `url(${searchSVG})` }"
                    :class="{ active: newSearchAvailable }" />
            </div>


        </div>

        <!-- Выбираем на какой конкретный аккаунт перевести -->
        <div class="account-id-block base-text" v-if="accounts && accounts.length">
            <span class="account-id-label small-text">account-id</span>
            <div class="select-wrapper">
                <span class="select-arrow"></span>
                <select class="account-select base-text" v-model="toAccountID">
                    <option v-for="account in accounts" :key="account.id">
                        {{ account.id }}
                    </option>
                </select>
            </div>
        </div>

        <!-- Указание суммы перевода в валюте аккаунта, с которого проводим -->
        <label class="amount-block base-text" v-if="accounts && accounts.length">
            amount
            <div class="amount-field">
                <div class="transfer-icon" :style="{ maskImage: `url(${transferSVG})` }" />
                <label class="currency base-text">
                    {{ currencySymbol }}
                    <input class="currency-input base-text" type="number" id="amount" placeholder="0.0"
                        v-model.number="amount" />
                </label>

            </div>
        </label>


        <span class="error base-text" v-if="errorText !== ''">
            {{ errorText }}
        </span>

        <div class="btns">
            <h3 class="cancel-btn btn" @click="$emit('cancel')"> cancel </h3>
            <h3 class="create-btn btn" @click="$emit('confirm', transfer)"> create </h3>
        </div>
    </div>
</template>

<script>
export default {
    name: 'NewTransferMenu',

    props: {
        accounts: {
            type: Array,
            required: false,
        },
        errorText: {
            type: String,
            required: false,
        },
    },

    data() {
        return {
            username: "",
            lastSearchedUsername: "",

            toAccountID: "",
            amount: null,

            searchSVG: require('@/assets/search.svg'),
            transferSVG: require('@/assets/transfer.svg'),

        }
    },

    computed: {
        newSearchAvailable() {
            return this.username !== this.lastSearchedUsername
        },
        currencySymbol() {
            const map = {
                USD: '$',
                EUR: '€',
                RUB: '₽'
            }
            if (!this.accounts || !this.accounts.length) return '';

            return map[this.accounts[0].currency] || this.accounts[0].currency;
        },
        transfer() {
            return {
                to_account_id: this.toAccountID,
                to_account_username: this.username,
                amount: this.amount,
                currency: this.accounts[0].currency,
            }
        }
    }
}
</script>

<style scoped>
.new-transfer {
    display: flex;
    flex-direction: column;
    flex: 1;
    max-width: fit-content;
    max-height: fit-content;
    align-content: center;

    gap: 5px;

    align-items: center;

    background-color: var(--color-bg-alt);
    border-radius: 10px;
    padding: 4px 10px;
}

.username-input-block {
    display: flex;
    flex: 1;
    flex-direction: row;
    align-self: stretch;
    justify-content: space-between;
    gap: 10px;

    background-color: var(--color-bg);
    border-radius: 10px;
    padding: 2px 5px;
}

.username-input-block .username-input {
    color: var(--color-text-muted);
    flex: 1;
}

.username-input-block .search-btn {
    flex: 1;
    max-width: fit-content;
    padding: 3px;
    border-radius: 7px;
    align-content: center;
    background-color: transparent;
}

.username-input-block .search-btn.active {
    background-color: var(--color-accent);
}

.username-input-block .search-icon {
    width: 14px;
    height: 14px;
    opacity: 1;

    pointer-events: none;
    transition: opacity 0.4s ease;
    background-color: var(--color-accent);
    mask-size: contain;
}

.username-input-block .search-icon.active {
    background-color: var(--color-bg);
}

.account-id-block {
    display: flex;
    flex: 1;
    max-width: fit-content;
    flex-direction: column;
}

.account-id-block .account-id-label {
    flex: 1;
    max-width: fit-content;
    text-align: left;
}

.select-wrapper {
    position: relative;
    width: fit-content;
}

.account-id-block .account-select {
    padding: 4px 6px 4px 18px;

    background-color: var(--color-bg);
    border-radius: 10px;

    border: none;
    outline: none;
    box-shadow: none;
    background-clip: padding-box;
    font: inherit;


    appearance: none;
    -webkit-appearance: none;
    -moz-appearance: none;
}

.account-id-block .select-arrow {
    position: absolute;
    top: 50%;
    left: 10px;
    transform: translate(-50%, -50%);


    width: 0px;
    height: 0px;
    pointer-events: none;

    border-left: 5px solid transparent;
    border-right: 5px solid transparent;
    border-top: 6px solid var(--color-accent);
}

.amount-block {
    display: flex;
    flex: 1;
    flex-direction: column;
    text-align: left;
    align-self: stretch;
}

.amount-block .amount-field {
    display: flex;
    flex: 1;
    flex-direction: row;
    align-self: stretch;

    align-items: center;

    background-color: var(--color-bg);
    border-radius: 10px;
    padding: 4px 10px;

    gap: 5px;

    appearance: none;
}

.amount-block .transfer-icon {
    width: 17px;
    height: 17px;
    opacity: 1;

    background-color: var(--color-accent);
    mask-size: contain;
}

.amount-block .currency {
    display: flex;
    flex: 1;
    align-self: stretch;
}

.amount-block .currency-input {
    flex: 1;
    align-self: stretch;
}

.error {
    flex: 1;
    align-self: stretch;
    text-align: right;
    color: var(--color-accent);
}

.btns {
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
    padding: 1px 7px;
}
</style>