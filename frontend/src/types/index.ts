export interface Project {
  id: string
  name: string
}

export interface Registry {
  id: string
  name: string
  status: string
  createdAt: string
  size: number
  // Extended for UI
  repositories?: Repository[]
  loadingRepos?: boolean
  expanded?: boolean
}

export interface Repository {
  name: string
  size: number
  updatedAt: string
}

export interface Image {
  digest: string
  tags: string[]
  size: number
  createdAt: string
}

export interface CleanupResult {
    deleted: unknown[]
    failed: unknown[]
}

export interface GCInfo {
  sizeNonReferenced: number
  sizeSummary: number
  sizeUntagged: number
}
