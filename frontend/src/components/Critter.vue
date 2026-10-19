<script setup lang="ts">
import type { CritterKind } from '@/scenes'

const VIEWBOX: Record<CritterKind, string> = {
  rabbit: '0 0 44 40',
  wolf: '0 0 80 50',
  frog: '0 0 32 24',
  butterfly: '0 0 30 24',
  firefly: '0 0 16 16',
  bird: '0 0 44 28',
  goat: '0 0 64 52',
  soldier: '0 0 34 58',
  jumper: '0 0 34 58',
  archer: '0 0 34 58',
  arrow: '0 0 30 8',
}

defineProps<{ kind: CritterKind }>()
</script>

<template>
  <div class="anchor">
    <div class="walker" :class="kind">
      <svg class="art" :class="kind" :viewBox="VIEWBOX[kind]">
        <template v-if="kind === 'rabbit'">
          <ellipse class="sh" cx="22" cy="38" rx="14" ry="2.4" />
          <circle cx="9" cy="27" r="5" fill="#ffffff" />
          <ellipse cx="20" cy="28" rx="13" ry="9" fill="#f4efe6" />
          <ellipse cx="20" cy="32" rx="11" ry="5" fill="#e0d8c8" />
          <circle cx="33" cy="20" r="8" fill="#f4efe6" />
          <ellipse cx="31" cy="7" rx="2.8" ry="8" fill="#ece6da" transform="rotate(-8 31 7)" />
          <ellipse cx="37" cy="8" rx="2.8" ry="8" fill="#f4efe6" transform="rotate(10 37 8)" />
          <ellipse cx="37" cy="8" rx="1.2" ry="5.5" fill="#f4a6b0" transform="rotate(10 37 8)" />
          <circle cx="36.5" cy="19" r="1.4" fill="#2a2a2a" />
          <circle cx="40.5" cy="22" r="1.3" fill="#f08a9a" />
          <ellipse cx="28" cy="36" rx="5" ry="2.6" fill="#e8e0d0" />
        </template>

        <template v-else-if="kind === 'wolf'">
          <ellipse class="sh" cx="40" cy="48" rx="28" ry="2.6" />
          <path d="M16,24 C6,16 3,26 1,31 C8,29 13,31 18,29 Z" fill="#59606e" />
          <rect class="leg a" x="20" y="32" width="5" height="15" rx="2" fill="#4b515d" />
          <rect class="leg b" x="46" y="32" width="5" height="15" rx="2" fill="#4b515d" />
          <ellipse cx="38" cy="27" rx="24" ry="11" fill="#6c7482" />
          <ellipse cx="38" cy="32" rx="20" ry="6" fill="#868e9c" />
          <path d="M22,20 L28,26 M30,17 L35,24 M38,16 L42,23" stroke="#5a616e" stroke-width="2" stroke-linecap="round" />
          <rect class="leg b" x="27" y="33" width="5" height="15" rx="2" fill="#6c7482" />
          <rect class="leg a" x="52" y="33" width="5" height="15" rx="2" fill="#6c7482" />
          <polygon points="54,14 74,23 77,29 60,33 52,27" fill="#6c7482" />
          <polygon points="56,15 60,6 64,17" fill="#59606e" />
          <polygon points="62,17 68,9 69,20" fill="#6c7482" />
          <polygon points="60,32 77,29 70,34" fill="#868e9c" />
          <circle cx="68" cy="23" r="1.7" fill="#f2d24a" />
          <circle cx="76.5" cy="27.5" r="1.6" fill="#1d2027" />
        </template>

        <template v-else-if="kind === 'frog'">
          <ellipse class="sh" cx="16" cy="22" rx="11" ry="2" />
          <ellipse cx="9" cy="19" rx="5" ry="3.5" fill="#4aa83c" />
          <ellipse cx="16" cy="15" rx="12" ry="7" fill="#5fbf4a" />
          <ellipse cx="16" cy="18" rx="8" ry="4" fill="#cfe9a0" />
          <circle cx="10" cy="9" r="4" fill="#5fbf4a" />
          <circle cx="22" cy="9" r="4" fill="#5fbf4a" />
          <circle cx="10" cy="8.5" r="2.5" fill="#fff" />
          <circle cx="22" cy="8.5" r="2.5" fill="#fff" />
          <circle cx="10.6" cy="8.5" r="1.2" fill="#222" />
          <circle cx="22.6" cy="8.5" r="1.2" fill="#222" />
          <path d="M9 15 Q16 19 23 15" fill="none" stroke="#2f7a2a" stroke-width="1.3" stroke-linecap="round" />
        </template>

        <template v-else-if="kind === 'butterfly'">
          <g class="wings">
            <ellipse cx="8" cy="9" rx="8" ry="7" fill="#ff9f43" />
            <ellipse cx="22" cy="9" rx="8" ry="7" fill="#ff9f43" />
            <ellipse cx="9" cy="17" rx="5.5" ry="4.5" fill="#ffd166" />
            <ellipse cx="21" cy="17" rx="5.5" ry="4.5" fill="#ffd166" />
            <circle cx="8" cy="9" r="2.6" fill="#7a3fb5" />
            <circle cx="22" cy="9" r="2.6" fill="#7a3fb5" />
          </g>
          <rect x="14" y="5" width="2" height="14" rx="1" fill="#3b2a1d" />
          <path d="M14.5 5 L11 1 M15.5 5 L19 1" stroke="#3b2a1d" stroke-width="1" stroke-linecap="round" />
        </template>

        <template v-else-if="kind === 'firefly'">
          <circle class="halo" cx="8" cy="8" r="7" fill="rgb(225 255 120 / 0.45)" />
          <circle cx="8" cy="8" r="2.4" fill="#f6ffa0" />
        </template>

        <template v-else-if="kind === 'bird'">
          <polygon points="12,15 2,9 4,20" fill="#2f3441" />
          <ellipse cx="22" cy="16" rx="10.5" ry="5.5" fill="#3d4352" />
          <ellipse cx="21" cy="18" rx="7" ry="3" fill="#8b93a6" />
          <circle cx="33" cy="12.5" r="4.4" fill="#3d4352" />
          <polygon points="36.5,11 43,13.5 36.5,15.5" fill="#f2b134" />
          <circle cx="34.5" cy="11.5" r="1.1" fill="#fff" />
          <path class="wing" d="M17,14 C13,0 27,-1 29,13 Z" fill="#2a2e3a" />
        </template>

        <template v-else-if="kind === 'goat'">
          <ellipse class="sh" cx="32" cy="50" rx="22" ry="2.4" />
          <rect class="leg a" x="14" y="34" width="4.5" height="15" rx="2" fill="#bdb6a6" />
          <rect class="leg b" x="38" y="34" width="4.5" height="15" rx="2" fill="#bdb6a6" />
          <ellipse cx="30" cy="30" rx="20" ry="10" fill="#efeadf" />
          <ellipse cx="30" cy="35" rx="16" ry="5" fill="#d9d2c2" />
          <rect class="leg b" x="21" y="35" width="4.5" height="14" rx="2" fill="#efeadf" />
          <rect class="leg a" x="44" y="35" width="4.5" height="14" rx="2" fill="#efeadf" />
          <path d="M10,25 C4,22 3,28 5,31 C8,29 10,30 12,30 Z" fill="#efeadf" />
          <path d="M46,24 L54,14 L58,20 L50,30 Z" fill="#efeadf" />
          <ellipse cx="55" cy="18" rx="8" ry="6" fill="#efeadf" />
          <ellipse cx="60" cy="21" rx="3.5" ry="2.6" fill="#d7cfbd" />
          <path d="M51,13 C47,5 53,0 58,3" fill="none" stroke="#a88e5a" stroke-width="2.6" stroke-linecap="round" />
          <path d="M55,12 C54,5 60,2 62,5" fill="none" stroke="#b79c6b" stroke-width="2.4" stroke-linecap="round" />
          <path d="M56,25 L57,30 L59,25 Z" fill="#d7cfbd" />
          <circle cx="56" cy="17" r="1.3" fill="#222" />
        </template>

        <template v-else-if="kind === 'soldier'">
          <ellipse class="sh" cx="17" cy="56" rx="9" ry="2" />
          <rect class="dark" x="11" y="35" width="4.4" height="20" rx="2" transform="rotate(8 13 35)" />
          <rect class="dark" x="17" y="35" width="4.4" height="20" rx="2" transform="rotate(-8 19 35)" />
          <path class="dark" d="M10,17 Q17,13 24,17 L22,37 L12,37 Z" />
          <rect class="dark" x="9.5" y="17" width="3.6" height="17" rx="1.8" transform="rotate(6 11 17)" />
          <circle class="dark" cx="17" cy="9" r="6.2" />
          <path class="dark-line" d="M26,34 L26,-4" />
          <polygon class="dark" points="24.4,-2 26,-8 27.6,-2" />
          <rect class="dark" x="22.5" y="17" width="3.6" height="17" rx="1.8" transform="rotate(-10 24 17)" />
          <g class="eyes">
            <ellipse cx="18.5" cy="8" rx="1.3" ry="1.6" fill="#ff2d2d" />
            <ellipse cx="22" cy="8" rx="1.3" ry="1.6" fill="#ff2d2d" />
          </g>
        </template>

        <template v-else-if="kind === 'jumper'">
          <ellipse class="sh" cx="17" cy="56" rx="9" ry="2" />
          <rect class="dark" x="11" y="35" width="4.4" height="20" rx="2" transform="rotate(14 13 35)" />
          <rect class="dark" x="17" y="35" width="4.4" height="20" rx="2" transform="rotate(-14 19 35)" />
          <path class="dark" d="M10,17 Q17,13 24,17 L22,37 L12,37 Z" />
          <rect class="dark" x="9.5" y="17" width="3.6" height="17" rx="1.8" transform="rotate(150 11 18)" />
          <circle class="dark" cx="17" cy="9" r="6.2" />
          <rect class="dark" x="22.5" y="17" width="3.6" height="17" rx="1.8" transform="rotate(-150 24 18)" />
          <g class="eyes">
            <ellipse cx="18.5" cy="8" rx="1.3" ry="1.6" fill="#ff2d2d" />
            <ellipse cx="22" cy="8" rx="1.3" ry="1.6" fill="#ff2d2d" />
          </g>
        </template>

        <template v-else-if="kind === 'archer'">
          <ellipse class="sh" cx="17" cy="56" rx="9" ry="2" />
          <rect class="dark" x="12" y="35" width="4.4" height="20" rx="2" />
          <rect class="dark" x="18" y="35" width="4.4" height="20" rx="2" />
          <path class="dark" d="M10,17 Q17,13 24,17 L22,37 L12,37 Z" />
          <circle class="dark" cx="17" cy="9" r="6.2" />
          <rect class="dark" x="22" y="17" width="3.6" height="17" rx="1.8" transform="rotate(-90 23.5 19)" />
          <rect class="dark" x="9.5" y="17" width="3.6" height="17" rx="1.8" transform="rotate(-95 11 19)" />
          <path class="dark-line" d="M32,8 C40,16 40,32 32,40" />
          <path class="dark-line thin" d="M32,8 L32,40" />
          <g class="eyes">
            <ellipse cx="18.5" cy="8" rx="1.3" ry="1.6" fill="#ff2d2d" />
            <ellipse cx="22" cy="8" rx="1.3" ry="1.6" fill="#ff2d2d" />
          </g>
        </template>

        <template v-else-if="kind === 'arrow'">
          <path class="dark-line" d="M2,4 H25" />
          <polygon class="dark" points="25,1.6 29,4 25,6.4" />
          <path class="dark-line thin" d="M2,4 L5,1.5 M2,4 L5,6.5" />
        </template>
      </svg>
    </div>
  </div>
</template>

<style scoped>
.anchor {
  position: absolute;
  transform: translate(-50%, -85%);
}
.walker {
  animation: wander var(--dur, 12s) linear var(--delay, 0s) infinite;
}
.walker.butterfly,
.walker.firefly {
  animation-name: drift;
  animation-timing-function: ease-in-out;
}
.art {
  display: block;
  width: 100%;
  height: auto;
  overflow: visible;
  transform-origin: 50% 100%;
}
.sh {
  fill: #000;
  opacity: 0.2;
}

.art.rabbit {
  animation: hop 0.8s ease-in-out infinite;
}
.art.wolf,
.art.goat {
  animation: trot 0.7s ease-in-out infinite;
}
.art.goat {
  animation-duration: 1s;
}
.art.frog {
  animation: frogHop 3.2s ease-in-out var(--delay, 0s) infinite;
}
.walker.arrow {
  animation: arrowX var(--dur, 3s) linear var(--delay, 0s) infinite;
}
.art.arrow {
  animation: arrowY var(--dur, 3s) ease-in-out var(--delay, 0s) infinite;
  transform-origin: 50% 50%;
}
.art.soldier {
  animation: stride 0.5s ease-in-out infinite;
}
.art.jumper {
  animation: jump 1.3s ease-in-out var(--delay, 0s) infinite;
}
.dark {
  fill: #0a0a12;
}
.dark-line {
  fill: none;
  stroke: #0a0a12;
  stroke-width: 1.7;
  stroke-linecap: round;
}
.dark-line.thin {
  stroke-width: 0.8;
}

.art.bird {
  animation: bob 0.9s ease-in-out infinite;
}
.leg {
  transform-box: fill-box;
  transform-origin: 50% 0%;
  animation: step 0.7s ease-in-out infinite alternate;
}
.leg.b {
  animation-direction: alternate-reverse;
}
.wings {
  transform-origin: 50% 50%;
  transform-box: view-box;
  animation: flutter 0.2s ease-in-out infinite alternate;
}
.wing {
  transform-origin: 22px 16px;
  animation: flap 0.3s ease-in-out infinite alternate;
}
.halo {
  animation: glow 2.4s ease-in-out var(--delay, 0s) infinite;
}
.art.firefly {
  filter: drop-shadow(0 0 4px rgb(225 255 120 / 0.9));
}
.eyes {
  transform-box: fill-box;
  transform-origin: 50% 50%;
  animation: blink 4.5s infinite;
  filter: drop-shadow(0 0 3px #ff2d2d);
}

@keyframes wander {
  0% {
    transform: translate(0, 0) scaleX(var(--face, 1));
  }
  46% {
    transform: translate(var(--dx), var(--dy)) scaleX(var(--face, 1));
  }
  50% {
    transform: translate(var(--dx), var(--dy)) scaleX(calc(var(--face, 1) * -1));
  }
  96% {
    transform: translate(0, 0) scaleX(calc(var(--face, 1) * -1));
  }
  100% {
    transform: translate(0, 0) scaleX(var(--face, 1));
  }
}
@keyframes arrowX {
  0% {
    transform: translateX(0) scaleX(var(--face, 1));
    opacity: 0;
  }
  4% {
    opacity: 1;
  }
  55% {
    transform: translateX(calc(var(--dx) * var(--face, 1))) scaleX(var(--face, 1));
    opacity: 1;
  }
  60%,
  100% {
    transform: translateX(calc(var(--dx) * var(--face, 1))) scaleX(var(--face, 1));
    opacity: 0;
  }
}
@keyframes arrowY {
  0% {
    transform: translateY(0) rotate(-28deg);
  }
  27% {
    transform: translateY(-46px) rotate(0deg);
  }
  55%,
  100% {
    transform: translateY(0) rotate(28deg);
  }
}
@keyframes stride {
  0%,
  100% {
    transform: translateY(0) rotate(-2deg);
  }
  50% {
    transform: translateY(-2px) rotate(2deg);
  }
}
@keyframes jump {
  0%,
  40% {
    transform: translateY(0) scale(1.1, 0.88);
  }
  55% {
    transform: translateY(-30px) scale(0.92, 1.1);
  }
  72% {
    transform: translateY(0) scale(1.12, 0.84);
  }
  86% {
    transform: translateY(0) scale(0.98, 1.04);
  }
  100% {
    transform: translateY(0) scale(1.1, 0.88);
  }
}
@keyframes drift {
  0%,
  100% {
    transform: translate(0, 0) rotate(0deg);
  }
  25% {
    transform: translate(var(--dx), calc(var(--dy) * -1)) rotate(14deg);
  }
  50% {
    transform: translate(calc(var(--dx) * 0.3), var(--dy)) rotate(-12deg);
  }
  75% {
    transform: translate(calc(var(--dx) * -0.6), calc(var(--dy) * 0.4)) rotate(8deg);
  }
}
@keyframes hop {
  0%,
  100% {
    transform: translateY(0) scale(1.1, 0.88);
  }
  50% {
    transform: translateY(-9px) scale(0.92, 1.1);
  }
}
@keyframes trot {
  0%,
  100% {
    transform: rotate(-3deg) scale(1, 1);
  }
  50% {
    transform: rotate(3deg) scale(1.03, 0.96);
  }
}
@keyframes step {
  from {
    transform: rotate(-24deg);
  }
  to {
    transform: rotate(24deg);
  }
}
@keyframes frogHop {
  0%,
  50%,
  100% {
    transform: translateY(0) scale(1, 1);
  }
  25% {
    transform: translateY(0) scale(1.04, 0.96);
  }
  60% {
    transform: translateY(0) scale(1.14, 0.8);
  }
  72% {
    transform: translateY(-15px) scale(0.88, 1.22);
  }
  84% {
    transform: translateY(0) scale(1.16, 0.78);
  }
  92% {
    transform: translateY(0) scale(0.98, 1.05);
  }
}
@keyframes bob {
  0%,
  100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-4px);
  }
}
@keyframes flutter {
  from {
    transform: scaleX(1);
  }
  to {
    transform: scaleX(0.25);
  }
}
@keyframes flap {
  from {
    transform: rotate(-35deg);
  }
  to {
    transform: rotate(28deg);
  }
}
@keyframes glow {
  0%,
  100% {
    opacity: 0.15;
  }
  50% {
    opacity: 1;
  }
}
@keyframes blink {
  0%,
  94%,
  100% {
    transform: scaleY(1);
  }
  97% {
    transform: scaleY(0.1);
  }
}
@media (prefers-reduced-motion: reduce) {
  .walker,
  .art,
  .leg,
  .wings,
  .wing,
  .halo,
  .eyes {
    animation: none;
  }
}
</style>
