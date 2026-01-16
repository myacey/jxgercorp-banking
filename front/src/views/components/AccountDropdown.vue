<template>
    <div class="account-dropdown btn" @click="$emit('click')">
        <span class="arrow" :class="{ open }">▶</span>
        <span class="base-text selector-text">account</span>
    </div>
</template>

<script>
export default {
    name: 'AccountDropdown',

    props: {
        open: {
            type: Boolean,
            default: false
        }
    }
};
</script>


<style scoped>
.account-dropdown {
    display: flex;
    flex-direction: row;
    gap: 2px;

    align-items: center;
    cursor: pointer;
}

.arrow {
    vertical-align: baseline;
    display: inline-block;
    transition: transform 0.25s ease;

    font-size: 10px;
    color: white;
}

.arrow.open {
    transform: rotate(180deg);
}

.selector-text {
    vertical-align: baseline;
}
</style>