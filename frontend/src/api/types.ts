import type { components } from './schema'

type Schemas = components['schemas']

export type GameMap = Schemas['GameMap']
export type MapId = Schemas['MapId']
export type MapNode = Schemas['MapNode']
export type Progress = Schemas['Progress']
export type Mode = Schemas['Mode']
export type StatsSummary = Schemas['StatsSummary']
export type Lang = Schemas['Lang']
export type CodeLanguage = Schemas['CodeLanguage']
export type TextItem = Schemas['TextItem']
export type SessionCreate = Schemas['SessionCreate']
export type SessionResult = Schemas['SessionResult']
export type SessionSummary = Schemas['SessionSummary']
export type TimelinePoint = Schemas['TimelinePoint']
export type KeyStat = Schemas['KeyStat']
export type Achievement = Schemas['Achievement']
export type Leaderboard = Schemas['Leaderboard']
export type User = Schemas['User']
export type StatsRange = 'week' | 'month' | 'year'

export interface TextQuery {
  mode: Mode
  lang: Lang
  difficulty: number
  count: number
  code_language?: CodeLanguage
}

export interface LeaderboardQuery {
  metric: 'speed' | 'accuracy' | 'xp'
  period: 'today' | 'week' | 'all'
  limit?: number
}

export interface Api {
  getMaps(): Promise<GameMap[]>
  getProgress(): Promise<Progress>
  getSummary(): Promise<StatsSummary>
  getSessions(limit?: number): Promise<SessionSummary[]>
  getTimeline(range: StatsRange): Promise<TimelinePoint[]>
  getKeys(kind: 'key' | 'bigram', limit?: number): Promise<KeyStat[]>
  getAchievements(): Promise<Achievement[]>
  getLeaderboard(query: LeaderboardQuery): Promise<Leaderboard>
  getMe(): Promise<User>
  updateMe(patch: { nickname: string }): Promise<User>
  reset?(): Promise<void>
  getTexts(query: TextQuery): Promise<TextItem[]>
  createSession(body: SessionCreate): Promise<SessionResult>
}
