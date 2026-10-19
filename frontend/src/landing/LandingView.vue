<script setup>
import AppHeader from './components/AppHeader.vue'
import HeroSection from './components/HeroSection.vue'
import FeaturesSection from './components/FeaturesSection.vue'
import GamificationSection from './components/GamificationSection.vue'
import AppFooter from './components/AppFooter.vue'
import './styles/landing.css'
</script>

<template>
  <div class="landing">
    <AppHeader />
    <main>
      <HeroSection />
      <FeaturesSection />
      <GamificationSection />
    </main>
    <AppFooter />
  </div>
</template>
