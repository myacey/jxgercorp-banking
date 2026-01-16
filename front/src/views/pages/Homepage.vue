<template>
  <div class="homepage">
    <!-- ОСНОВНОЕ ОКНО НАВИГАЦИИ -->
    <AppNavigation />
    <div class="divider">
      <div class="divider-line"></div>
    </div>


    <!-- ОСНОВНОЕ ОКНО ПРИЛОЖЕНИЯ -->
    <main class="main-content">

      <!-- ВЕРХНЕЕ ВСПОМОГАТЕЛЬНОЕ МЕНЮ-->
      <div class="top-selector">
        <!-- ОКНО ВЫБОРА АККАУНТА -->
        <div class="account-dropdown-wrapper">
          <AccountDropdown :open="showAccountSelector" @click="showAccountSelector = !showAccountSelector" />
          <AccountSelector v-if="showAccountSelector" :accounts="accounts" :selectedAccountID="selectedAccount?.id"
            @select="onSelectAccount" @account-created="addAccount" @request-delete="onDeleteAccountModal" />
        </div>
      </div>

      <section class="dashboard-section">
        <template v-if="selectedAccount">
          <NewTransferBlock :account="selectedAccount" @new-transfer="showNewTransferMenu = true" />
          <TransactionHistory :transactions="transactions" :accountID="selectedAccount.id" />
        </template>

        <template v-else>
          <p class="empty-account-state">
            no accounts yet. create a new one
          </p>
        </template>
      </section>
    </main>

    <!-- MODAL -->
    <!-- delete account -->
    <div v-if="showAccountDeleteConfirm" class="overlay" @click.self="showAccountDeleteConfirm = false">
      <DeleteAccountModal :account="accountToDelete" @cancel="showAccountDeleteConfirm = false"
        @confirm="onDeleteAccount" />
    </div>

    <!-- new transfer -->
    <div v-if="showNewTransferMenu" class="overlay"
      @click.self="showNewTransferMenu = false; accountsToTransfer = null">
      <NewTransferMenu :accounts=accountsToTransfer :errorText=errorOnCreateTransfer
        @request-search-accounts="searchAccountsOfUsername"
        @cancel="showNewTransferMenu = false; accountsToTransfer = null" @confirm="onCreateTransfer" />
    </div>

  </div>
</template>

<script>
import AppNavigation from '../components/AppNavigation.vue'
import TransactionHistory from '../components/TransactionHistory.vue'
import NewTransferBlock from '../components/NewTransferBlock.vue'
import AccountDropdown from '../components/AccountDropdown.vue';
import AccountSelector from '../components/AccountSelector.vue';
import { createAccount, createTransfer, deleteAccount, fetchAccounts, fetchTransfers } from '@/api/api';
import { getUsernameFromToken } from '@/utils/auth';
import DeleteAccountModal from '../components/DeleteAccountModal.vue';
import NewTransferMenu from '../components/NewTransferMenu.vue';
import { convertKeysToCamel } from '@/utils/snake2camel';

export default {
  name: 'HomePage',

  components: {
    AppNavigation,
    NewTransferBlock,
    TransactionHistory,
    AccountDropdown,
    AccountSelector,
    DeleteAccountModal,
    NewTransferMenu
  },

  data() {
    return {
      accounts: [],
      selectedAccount: null,
      showAccountSelector: false,
      showAccountDeleteConfirm: false,
      accountToDelete: null,


      // ------------ ПЕРЕВОДЫ ------------
      transactions: [],

      // ------ создание нового перевода ------
      showNewTransferMenu: false,
      // список аккаунтов искомого пользователя, на который можно перевести деньги
      accountsToTransfer: [],
      errorOnCreateTransfer: "",
    }
  },

  async mounted() {
    await this.loadAccounts()
  },

  watch: {
    selectedAccount: function (newAccount) {
      if (newAccount)
        this.onFetchTransactions()
      else
        this.transactions = []
    },
    immediate: true
  },

  methods: {
    // ------------ АККАУНТЫ ------------
    async loadAccounts() {
      try {
        const params = {
          username: getUsernameFromToken()
        }
        const resp = await fetchAccounts(params)
        this.accounts = resp

        if (this.accounts.length != 0) {
          this.selectedAccount = this.accounts[0]
        }
        console.log('SELECTED ACCOUNT:', this.selectedAccount)
      } catch (error) {
        console.error("failed to get accounts")
      }
    },
    onSelectAccount(account) {
      this.selectedAccount = account
      this.showAccountSelector = false
    },
    async addAccount({ currency_code }) {
      try {
        const data = {
          currency: currency_code,
        }
        const newAccount = await createAccount(data)
        this.accounts.push(newAccount)
        this.selectedAccount = newAccount
      } catch (error) {
        console.error("failed to create account")
      }
    },
    onDeleteAccountModal(account) {
      this.showAccountDeleteConfirm = true;
      this.accountToDelete = account;
    },
    async onDeleteAccount(account) {
      try {
        const data = {
          account_id: account.id
        }
        await deleteAccount(data);

        // чистим список аккаунтов
        this.accounts = this.accounts.filter(
          acc => acc.id !== account.id
        )

        // если удаленный аккаунт - выбранный
        if (this.selectedAccount?.id === account.id) {
          this.selectedAccount = this.accounts.length > 0
            ? this.accounts[0]
            : null
        }

        // закрытие менюшек и сброс всего
        this.showAccountDeleteConfirm = false
        this.accountToDelete = null

      } catch (error) {
        console.error('failed to delete account:', error)
      }
    },
    async searchAccountsOfUsername(username) {
      this.errorOnCreateTransfer = ''
      try {
        const params = {
          username: username,
          currency: this.selectedAccount.currency,
        }
        this.accountsToTransfer = await fetchAccounts(params)

        if (this.accountsToTransfer.length == 0) {
          this.errorOnCreateTransfer = "no accounts found"
        }
      } catch (error) {
        this.errorOnCreateTransfer = error.data?.message;
        console.error('failed to fetch accounts to transfer:', error)
      }
      console.log('aaaa:', this.errorOnCreateTransfer)
    },
    // ------------ ПЕРЕВОДЫ ------------
    async onFetchTransactions() {
      try {
        const params = {
          current_account_id: this.selectedAccount.id,
          offset: 0,
          limit: 20,
        }

        const resp = await fetchTransfers(params)
        this.transactions = convertKeysToCamel(resp)
      } catch (error) {
        console.error("failed to fetch transactions:", error)
      }
    },
    async onCreateTransfer(transfer) {
      this.errorOnCreateTransfer = ''
      try {
        const data = {
          from_account_id: this.selectedAccount.id,
          from_account_username: this.selectedAccount.owner_username,
          ...transfer,
        }

        const resp = await createTransfer(data)
        this.showNewTransferMenu = false
        this.selectedAccount.balance -= resp.amount

        this.onFetchTransactions()
      } catch (error) {
        this.errorOnCreateTransfer = error.data?.message;
        console.error('failed to create transfer:', error)
      }
    }
  }
}
</script>

<style scoped>
.homepage {
  min-height: 100vh;
  background-color: rgba(40, 40, 40, 1);
  display: flex;
  flex-direction: column;
}

.divider {
  display: flex;
  margin-top: 10px;
  width: 100%;
  flex-direction: column;
  overflow: hidden;
  align-items: stretch;
  justify-content: center;
}

@media (max-width: 991px) {
  .divider {
    max-width: 100%;
  }
}

.divider-line {
  border-radius: 5px;
  background-color: rgba(217, 217, 217, 1);
  display: flex;
  min-height: 1px;
  width: 100%;
}

@media (max-width: 991px) {
  .divider-line {
    max-width: 100%;
  }
}

.main-content {
  display: inline-flex;
  flex-direction: column;

  border-radius: 10px;
  background-color: var(--color-bg-alt);
  margin: auto;
  padding: 5px;
  max-width: fit-content;
}

.top-selector {
  display: flex;
  flex: 1;
  margin: 5px;
  align-self: stretch;
}

.dashboard-section {
  align-items: center;
  display: flex;
  gap: 40px 98px;
  justify-content: start;
  flex-wrap: wrap;
  margin: 15px;
}

.account-dropdown-wrapper {
  position: relative;
}

.overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1001;
}
</style>
