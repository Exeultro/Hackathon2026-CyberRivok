import { createRouter, createWebHistory } from 'vue-router'
import LandingView from '@/landing/LandingView.vue'
import MapView from '@/views/MapView.vue'
import PlayView from '@/views/PlayView.vue'
import PracticeView from '@/views/PracticeView.vue'
import StatsView from '@/views/StatsView.vue'
import LeaderboardView from '@/views/LeaderboardView.vue'
import ProfileView from '@/views/ProfileView.vue'
import SettingsView from '@/views/SettingsView.vue'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'landing', component: LandingView, meta: { bare: true } },
    { path: '/map', name: 'map', component: MapView },
    { path: '/play/:nodeId', name: 'play', component: PlayView, props: true },
    { path: '/practice', name: 'practice', component: PracticeView },
    { path: '/stats', name: 'stats', component: StatsView },
    { path: '/leaderboard', name: 'leaderboard', component: LeaderboardView },
    { path: '/profile', name: 'profile', component: ProfileView },
    { path: '/settings', name: 'settings', component: SettingsView },
  ],
})
