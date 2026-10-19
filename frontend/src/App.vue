<script setup lang="ts">
import { useRoute } from 'vue-router'
import AppHeader from '@/components/AppHeader.vue'
import SideNav from '@/components/SideNav.vue'

const route = useRoute()
</script>

<template>
  <RouterView v-if="route.meta.bare" />
  <template v-else-if="route.matched.length">
    <AppHeader />
    <SideNav />
    <div class="content">
      <RouterView />
    </div>
  </template>
</template>

<style scoped>
.content {
  margin-left: var(--nav-w);
  padding-top: var(--head-h);
  padding-bottom: var(--nav-b);
  min-height: 100vh;
}
</style>
