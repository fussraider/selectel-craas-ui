import type { Repository } from './index'

export interface RepoNode {
  name: string
  fullPath: string
  children: RepoNode[]
  repo?: Repository
  isGroup: boolean
}
