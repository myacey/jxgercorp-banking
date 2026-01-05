<template>
  <div class="homepage">
    <!-- ОСНОВНОЕ ОКНО НАВИГАЦИИ -->
    <AppNavigation />
    <div class="divider">
      <div class="divider-line"></div>
    </div>

    <!-- ОКНО ВЫБОРА АККАУНТА -->

    <!-- ОСНОВНОЕ ОКНО ПРИЛОЖЕНИЯ -->
    <main class="main-content">
      <div class="top-selector">
        <div class="account-dropdown-wrapper">
          <AccountDropdown @click="showAccountSelector = !showAccountSelector" />
          <AccountSelector v-if="showAccountSelector" :accounts="accounts" :selectedAccountID="selectedAccount.id"
            @select="onSelectAccount" />
        </div>
      </div>

      <section class="dashboard-section">
        <NewTransferBlock v-if="selectedAccount" :account="selectedAccount" />
        <TransactionHistory v-if="selectedAccount" :accountID="selectedAccount.id" />
      </section>
    </main>

  </div>
</template>

<script>
import AppNavigation from '../components/AppNavigation.vue'
import TransactionHistory from '../components/TransactionHistory.vue'
import NewTransferBlock from '../components/NewTransferBlock.vue'
import AccountDropdown from '../components/AccountDropdown.vue';
import AccountSelector from '../components/AccountSelector.vue';
import { fetchAccounts } from '@/api/api';
import { getUsernameFromToken } from '@/utils/auth';

export default {
  name: 'HomePage',

  components: {
    AppNavigation,
    NewTransferBlock,
    TransactionHistory,
    AccountDropdown,
    AccountSelector
  },

  data() {
    return {
      accounts: [],
      selectedAccount: null,
      showAccountSelector: false
    }
  },

  mounted() {
    this.loadAccounts()
  },

  methods: {
    async loadAccounts() {
      try {
        const params = {
          username: getUsernameFromToken()
        }
        const resp = await fetchAccounts(params)
        this.accounts = resp

        if (this.accounts.length > 1) {
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
</style>
