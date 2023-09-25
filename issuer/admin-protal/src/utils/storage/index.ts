import { createStorage } from "./storage"
export * from "./types"

export const storage = createStorage(localStorage)

export const session = createStorage(sessionStorage)
