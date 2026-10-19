import { httpApi } from './httpApi'
import { mockApi } from './mockApi'
import type { Api } from './types'

// VITE_API_URL empty -> localStorage mock (VITE_API_URL= npm run dev)
export const api: Api = import.meta.env.VITE_API_URL ? httpApi : mockApi
export type {
  Achievement,
  Api,
  CodeLanguage,
  GameMap,
  KeyStat,
  Lang,
  Leaderboard,
  LeaderboardQuery,
  MapId,
  MapNode,
  Mode,
  Progress,
  SessionCreate,
  SessionResult,
  SessionSummary,
  StatsRange,
  StatsSummary,
  TextItem,
  TimelinePoint,
  User,
} from './types'
