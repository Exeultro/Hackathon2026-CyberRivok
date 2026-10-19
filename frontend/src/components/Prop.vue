<script setup lang="ts">
import type { PropKind } from '@/scenes'

const VIEWBOX: Record<PropKind, string> = {
  oak: '0 0 60 64',
  pine: '0 0 60 80',
  bush: '0 0 50 36',
  flowers: '0 0 40 24',
  tuft: '0 0 24 16',
  mushroom: '0 0 24 24',
  stump: '0 0 30 26',
  log: '0 0 50 22',
  pond: '0 0 70 40',
  fence: '0 0 60 24',
  signpost: '0 0 30 44',
  hill: '0 0 140 70',
  haystack: '0 0 50 44',
  rock: '0 0 40 32',
  cairn: '0 0 24 28',
  peak: '0 0 100 100',
  torch: '0 0 20 40',
  banner: '0 0 30 56',
  tent: '0 0 60 50',
  campfire: '0 0 30 28',
  barrel: '0 0 20 26',
  catapult: '0 0 70 50',
}

withDefaults(defineProps<{ kind: PropKind; variant?: number; sway?: boolean }>(), { variant: 0, sway: false })
</script>

<template>
  <svg
    class="prop"
    :class="{ sway, grass: sway && (kind === 'tuft' || kind === 'flowers') }"
    :viewBox="VIEWBOX[kind]"
  >
    <template v-if="kind === 'oak'">
      <ellipse class="shadow" cx="30" cy="54" rx="22" ry="8" />
      <g>
      <rect class="trunk" x="26" y="40" width="8" height="16" rx="2" />
      <g class="flip">
        <circle class="c f3" cx="30" cy="28" r="24" />
        <circle class="c f3" cx="12" cy="36" r="11" />
        <circle class="c f3" cx="48" cy="36" r="11" />
        <circle class="c f3" cx="16" cy="17" r="10" />
        <circle class="c f3" cx="44" cy="17" r="10" />
        <circle class="c f2" cx="30" cy="26" r="20" />
        <circle class="c f1" cx="24" cy="20" r="12" />
        <circle class="c f1 soft" cx="19" cy="15" r="4" />
      </g>
      </g>
    </template>

    <template v-else-if="kind === 'pine'">
      <ellipse class="shadow" cx="30" cy="72" rx="16" ry="6" />
      <g>
      <rect class="trunk" x="27" y="62" width="6" height="12" rx="1" />
      <g class="flip">
        <polygon class="c f2" points="30,26 4,66 30,66" />
        <polygon class="c f3" points="30,26 56,66 30,66" />
        <polygon class="c f1" points="30,14 8,48 30,48" />
        <polygon class="c f3" points="30,14 52,48 30,48" />
        <polygon class="c f1" points="30,2 13,30 30,30" />
        <polygon class="c f2" points="30,2 47,30 30,30" />
      </g>
      </g>
    </template>

    <template v-else-if="kind === 'bush'">
      <ellipse class="shadow" cx="25" cy="30" rx="19" ry="5" />
      <g>
      <g class="flip">
        <circle class="c f3" cx="14" cy="22" r="12" />
        <circle class="c f3" cx="36" cy="22" r="12" />
        <circle class="c f2" cx="25" cy="16" r="14" />
        <circle class="c f1" cx="20" cy="12" r="8" />
      </g>
      </g>
    </template>

    <template v-else-if="kind === 'flowers'">
      <g>
        <path class="stroke f2" d="M6 22 L8 14 M18 22 L20 10 M30 22 L31 17" />
        <circle cx="8" cy="13" r="3.2" fill="#f472b6" />
        <circle cx="20" cy="8" r="3.2" fill="#fde047" />
        <circle cx="31" cy="16" r="3.2" fill="#ffffff" />
      </g>
    </template>

    <template v-else-if="kind === 'tuft'">
      <g>
        <path class="stroke f2" d="M4 14 L6 4 M12 14 L12 2 M20 14 L18 5" />
      </g>
    </template>

    <template v-else-if="kind === 'mushroom'">
      <ellipse class="shadow" cx="12" cy="21" rx="8" ry="2.5" />
      <rect x="9.5" y="12" width="5" height="9" rx="2" fill="#f3e9d2" />
      <path d="M3,13 C3,3 21,3 21,13 Z" fill="#e5484d" />
      <circle cx="8" cy="9" r="1.6" fill="#fff" />
      <circle cx="14" cy="7.5" r="1.4" fill="#fff" />
      <circle cx="17.5" cy="11" r="1.2" fill="#fff" />
    </template>

    <template v-else-if="kind === 'stump'">
      <ellipse class="shadow" cx="15" cy="22" rx="12" ry="3.5" />
      <rect class="trunk" x="5" y="10" width="20" height="12" rx="3" />
      <ellipse cx="15" cy="10" rx="10" ry="5" fill="#d3a56a" />
      <ellipse cx="15" cy="10" rx="5.5" ry="2.6" fill="none" stroke="#a57a42" stroke-width="1.2" />
    </template>

    <template v-else-if="kind === 'log'">
      <ellipse class="shadow" cx="25" cy="19" rx="22" ry="3.5" />
      <rect class="trunk" x="4" y="6" width="40" height="13" rx="6" />
      <ellipse cx="44" cy="12.5" rx="4" ry="6.5" fill="#d3a56a" />
      <ellipse cx="44" cy="12.5" rx="1.8" ry="3" fill="none" stroke="#a57a42" stroke-width="1" />
      <path d="M12 9 L26 9 M16 14 L34 14" stroke="#5e3b1c" stroke-width="1.2" opacity="0.5" />
    </template>

    <template v-else-if="kind === 'pond'">
      <ellipse cx="35" cy="22" rx="33" ry="17" fill="#cdb27a" />
      <ellipse cx="35" cy="22" rx="30" ry="14.5" fill="#3aa3da" />
      <ellipse cx="30" cy="19" rx="21" ry="9" fill="#62c6f0" />
      <path d="M20 22 Q26 20 32 22 M38 26 Q44 24 50 26" fill="none" stroke="#fff" stroke-width="1.4" opacity="0.7" stroke-linecap="round" />
      <path class="stroke f2" d="M60 14 L61 6 M63 16 L66 8" />
    </template>

    <template v-else-if="kind === 'fence'">
      <ellipse class="shadow" cx="30" cy="22" rx="27" ry="2.5" />
      <rect x="4" y="9" width="52" height="3" fill="#d8ad6e" />
      <rect x="4" y="15" width="52" height="3" fill="#d8ad6e" />
      <rect x="6" y="5" width="5" height="17" rx="1" fill="#b98a52" />
      <rect x="27" y="5" width="5" height="17" rx="1" fill="#b98a52" />
      <rect x="49" y="5" width="5" height="17" rx="1" fill="#b98a52" />
    </template>

    <template v-else-if="kind === 'signpost'">
      <ellipse class="shadow" cx="15" cy="41" rx="7" ry="2.5" />
      <rect class="trunk" x="13" y="10" width="4" height="31" />
      <polygon points="3,6 24,6 29,12 24,18 3,18" fill="#d8ad6e" stroke="#8a5a2b" stroke-width="1.2" />
      <path d="M8 12 H20" stroke="#8a5a2b" stroke-width="1.6" stroke-linecap="round" />
    </template>

    <template v-else-if="kind === 'hill' && variant === 0">
      <path class="c f3" d="M2,62 C8,24 44,6 70,6 C100,6 132,24 138,62 Q70,78 2,62 Z" />
      <g class="flip">
        <path class="c f2" d="M2,62 C8,24 44,6 70,6 C54,22 42,44 40,69 Q18,68 2,62 Z" />
        <path class="c f1" d="M16,52 C22,32 42,16 64,10 C50,24 38,38 34,56 Q24,55 16,52 Z" />
      </g>
      <path class="rim" d="M2,62 Q70,78 138,62" />
      <circle cx="94" cy="32" r="2.2" fill="#fff" />
      <circle cx="106" cy="44" r="2.2" fill="#f9a8d4" />
      <circle cx="84" cy="46" r="2.2" fill="#fde047" />
      <path class="stroke f1" d="M72 28 L73 21 M78 32 L80 25" />
    </template>

    <template v-else-if="kind === 'hill' && variant === 1">
      <path class="c f3" d="M2,62 C4,40 26,28 46,36 C56,14 92,12 104,34 C122,28 136,42 138,62 Q70,76 2,62 Z" />
      <g class="flip">
        <path class="c f2" d="M2,62 C4,40 26,28 46,36 C40,46 36,56 36,66 Q18,66 2,62 Z" />
        <path class="c f2" d="M46,36 C56,14 92,12 104,34 C90,26 72,24 62,32 C54,40 52,52 52,68 L44,67 Z" />
        <path class="c f1" d="M12,54 C16,44 30,38 40,40 C30,46 24,52 22,58 Q16,57 12,54 Z" />
      </g>
      <path class="rim" d="M2,62 Q70,76 138,62" />
      <circle cx="116" cy="46" r="2.2" fill="#fff" />
      <circle cx="124" cy="54" r="2.2" fill="#f9a8d4" />
      <path class="stroke f1" d="M80 20 L82 13 M88 22 L90 15" />
    </template>

    <template v-else-if="kind === 'hill'">
      <path class="c f3" d="M4,62 C12,30 46,10 74,10 C104,10 130,30 136,62 Q70,76 4,62 Z" />
      <g class="flip">
        <path class="c f2" d="M4,62 C12,30 46,10 74,10 C58,26 46,46 44,68 Q22,67 4,62 Z" />
        <path class="c f1" d="M18,52 C26,36 46,22 66,16 C52,30 42,44 38,56 Q28,55 18,52 Z" />
      </g>
      <polygon class="c r2" points="98,40 106,32 118,38 122,50 100,52" />
      <polygon class="c r1" points="106,32 118,38 112,42 104,40" />
      <polygon class="c f3" points="64,2 56,18 72,18" />
      <polygon class="c f2" points="64,-6 58,6 70,6" />
      <rect class="trunk" x="62" y="16" width="4" height="5" />
      <path class="rim" d="M4,62 Q70,76 136,62" />
    </template>

    <template v-else-if="kind === 'haystack'">
      <ellipse class="shadow" cx="25" cy="38" rx="20" ry="5" />
      <path class="c y2" d="M5,36 C5,8 45,8 45,36 Z" />
      <g class="flip">
        <path class="c y1" d="M5,36 C5,8 26,8 26,10 C14,16 12,28 14,36 Z" />
      </g>
      <path class="stroke y2" d="M14 30 L36 26 M12 22 L34 18" />
    </template>

    <template v-else-if="kind === 'rock'">
      <ellipse class="shadow" cx="20" cy="28" rx="16" ry="4" />
      <polygon class="c r3" points="4,26 8,12 20,6 34,12 37,26" />
      <g class="flip">
        <polygon class="c r2" points="8,12 20,6 24,16 14,24 4,26" />
        <polygon class="c r1" points="10,12 20,8 22,13 14,16" />
      </g>
    </template>

    <template v-else-if="kind === 'cairn'">
      <ellipse class="shadow" cx="12" cy="25" rx="10" ry="3" />
      <ellipse class="c r3" cx="12" cy="21" rx="10" ry="5.5" />
      <ellipse class="c r2" cx="12" cy="14.5" rx="7.5" ry="4.5" />
      <ellipse class="c r1" cx="11" cy="8.5" rx="5" ry="3.5" />
    </template>

    <template v-else-if="kind === 'peak' && variant === 0">
      <ellipse class="shadow" cx="50" cy="94" rx="48" ry="7" />
      <ellipse class="c r3" cx="50" cy="91" rx="49" ry="9" />
      <g class="flip">
        <polygon class="c r2" points="50,4 6,90 50,99" />
        <polygon class="c r3" points="50,4 50,99 94,90" />
        <polygon class="c r1" points="50,4 24,54 38,62 50,26" />
        <polygon class="c s1" points="50,4 33,38 41,33 46,42 50,32" />
        <polygon class="c s2" points="50,4 50,32 55,41 60,33 67,38" />
      </g>
    </template>

    <template v-else-if="kind === 'peak' && variant === 1">
      <ellipse class="shadow" cx="50" cy="94" rx="48" ry="7" />
      <ellipse class="c r3" cx="50" cy="91" rx="49" ry="9" />
      <g class="flip">
        <polygon class="c r2" points="74,34 48,92 74,99" />
        <polygon class="c r3" points="74,34 74,99 98,90" />
        <polygon class="c s1" points="74,34 66,50 71,47 74,53" />
        <polygon class="c r2" points="38,6 2,90 38,99" />
        <polygon class="c r3" points="38,6 38,99 76,90" />
        <polygon class="c r1" points="38,6 16,56 28,62 38,28" />
        <polygon class="c s1" points="38,6 24,36 31,32 35,40 38,31" />
        <polygon class="c s2" points="38,6 38,31 43,39 48,32 53,36" />
      </g>
    </template>

    <template v-else-if="kind === 'torch'">
      <ellipse class="shadow" cx="10" cy="36" rx="6" ry="2.5" />
      <rect class="trunk" x="8.5" y="16" width="3" height="20" />
      <path d="M10,2 C16,8 15,14 10,16 C5,14 4,8 10,2 Z" fill="#ffb347" />
      <path d="M10,7 C13,10 12,13 10,14 C8,13 7,10 10,7 Z" fill="#fff3b0" />
    </template>

    <template v-else-if="kind === 'banner'">
      <ellipse class="shadow" cx="15" cy="52" rx="7" ry="2.5" />
      <rect class="trunk" x="14" y="2" width="3" height="50" />
      <path d="M17,6 C22,4 24,9 29,7 L29,22 C24,24 22,19 17,21 Z" fill="#c62828" />
    </template>

    <template v-else-if="kind === 'tent'">
      <ellipse class="shadow" cx="30" cy="45" rx="26" ry="5" />
      <polygon points="30,6 3,46 30,48" fill="#e3d4ac" />
      <polygon points="30,6 30,48 57,46" fill="#bfae82" />
      <polygon points="30,24 21,47 39,47" fill="#3b3226" />
      <path d="M30 6 V1" stroke="#7a4f2a" stroke-width="1.5" />
      <path d="M30 1 L38 3 L30 5 Z" fill="#c62828" />
    </template>

    <template v-else-if="kind === 'campfire'">
      <ellipse class="shadow" cx="15" cy="25" rx="11" ry="3" />
      <circle class="c r2" cx="6" cy="21" r="3.2" />
      <circle class="c r2" cx="11" cy="24" r="3.2" />
      <circle class="c r2" cx="19" cy="24" r="3.2" />
      <circle class="c r2" cx="24" cy="21" r="3.2" />
      <path d="M8 22 L22 18 M8 18 L22 22" stroke="#5e3b1c" stroke-width="3" stroke-linecap="round" />
      <path d="M15,2 C23,10 22,18 15,20 C8,18 7,10 15,2 Z" fill="#ff9a3c" />
      <path d="M15,9 C19,13 18,17 15,18 C12,17 11,13 15,9 Z" fill="#fff0a8" />
    </template>

    <template v-else-if="kind === 'barrel'">
      <ellipse class="shadow" cx="10" cy="23" rx="8" ry="2.5" />
      <rect x="3" y="6" width="14" height="16" rx="3" fill="#8a5a2b" />
      <rect x="3" y="10" width="14" height="2" fill="#4f3218" />
      <rect x="3" y="17" width="14" height="2" fill="#4f3218" />
      <ellipse cx="10" cy="6" rx="7" ry="2.6" fill="#a8743a" />
    </template>

    <template v-else-if="kind === 'catapult'">
      <ellipse class="shadow" cx="35" cy="45" rx="30" ry="4" />
      <rect class="trunk" x="8" y="34" width="54" height="6" rx="2" />
      <polygon points="26,34 40,34 34,16" fill="#8a5a2b" />
      <path d="M24 34 L54 8" stroke="#a8743a" stroke-width="5" stroke-linecap="round" />
      <ellipse cx="55" cy="8" rx="7" ry="4" fill="#5a3a1c" />
      <circle cx="16" cy="41" r="6" fill="#5a3a1c" />
      <circle cx="16" cy="41" r="2.4" fill="#a8743a" />
      <circle cx="54" cy="41" r="6" fill="#5a3a1c" />
      <circle cx="54" cy="41" r="2.4" fill="#a8743a" />
    </template>
  </svg>
</template>

<style scoped>
.prop {
  position: absolute;
  height: auto;
  overflow: visible;
  transform: translate(-50%, -85%) scaleY(var(--sy, 1));
  transform-origin: 50% 85%;
}
.shadow {
  fill: #000;
  opacity: var(--sho, 0.25);
  transform: translate(calc(var(--shx, 0) * 0.3px), 0);
}
.prop.sway {
  animation: sway calc(2.6s + var(--tone, 0.5) * 2.4s) ease-in-out calc(var(--tone, 0.5) * -6s) infinite alternate;
}
.prop.sway.grass {
  animation-name: sway-grass;
}
@keyframes sway {
  from {
    rotate: -0.8deg;
  }
  to {
    rotate: 2.2deg;
  }
}
@keyframes sway-grass {
  from {
    rotate: -5deg;
  }
  to {
    rotate: 9deg;
  }
}
@media (prefers-reduced-motion: reduce) {
  .prop.sway {
    animation: none;
  }
}
.flip {
  transform: scaleX(var(--flip, 1));
  transform-origin: 50% 50%;
  transform-box: view-box;
}
.trunk {
  fill: color-mix(in srgb, var(--edge) calc(var(--depth, 0) * 40%), var(--trunk));
}
.c {
  fill: color-mix(
    in srgb,
    var(--edge) calc(var(--depth, 0) * 40%),
    color-mix(in srgb, var(--base) calc(90% + var(--tone, 0.5) * 10%), #fff)
  );
}
.rim {
  fill: none;
  stroke: var(--edge);
  stroke-width: 2;
  stroke-linecap: round;
  opacity: 0.3;
}
.stroke {
  fill: none;
  stroke: var(--base);
  stroke-width: 2;
  stroke-linecap: round;
}
.line {
  fill: none;
  stroke: var(--base);
  stroke-width: 1.2;
  opacity: 0.7;
}
.soft {
  opacity: 0.5;
}
.f1 { --base: var(--f1); }
.f2 { --base: var(--f2); }
.f3 { --base: var(--f3); }
.r1 { --base: var(--r1); }
.r2 { --base: var(--r2); }
.r3 { --base: var(--r3); }
.s1 { --base: var(--s1); }
.s2 { --base: var(--s2); }
.w1 { --base: var(--w1); }
.w2 { --base: var(--w2); }
.w3 { --base: var(--w3); }
.y1 { --base: var(--y1); }
.y2 { --base: var(--y2); }
</style>
