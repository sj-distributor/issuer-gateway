export interface Updater<T> {
  (previousState?: T): T
}

export interface Options<T> {
  serializer?: (value: T) => string
  deserializer?: (value: string) => T
  defaultValue?: T | Updater<T>
}

export enum StorageKeys {
  TOKEN = "TOKEN",
}

export type StorageKeyType = StorageKeys.TOKEN
