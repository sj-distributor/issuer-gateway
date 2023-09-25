import { Options, StorageKeyType } from "./types"

// eslint-disable-next-line @typescript-eslint/ban-types
function isFunction(value: unknown): value is Function {
  return typeof value === "function"
}

function serializer<T>(value: T, options?: Options<T>) {
  if (options?.serializer) {
    return options?.serializer(value)
  }
  return JSON.stringify(value)
}

function deserializer<T>(value: string, options?: Options<T>): T {
  if (options?.deserializer) {
    return options.deserializer(value)
  }
  return JSON.parse(value)
}

export const createStorage = (storage: Storage = localStorage) => {
  function get<T>(key: StorageKeyType, options?: Options<T>): T | null {
    try {
      const raw = storage.getItem(key)
      if (raw) {
        return deserializer(raw, options)
      }
    } catch (e) {
      console.error(e)
    }
    if (isFunction(options?.defaultValue)) {
      return options?.defaultValue() ?? null
    }
    return options?.defaultValue ?? null
  }

  function set<T>(key: StorageKeyType, value: T, options?: Options<T>) {
    try {
      storage.setItem(key, serializer(value, options))
    } catch (e) {
      console.error(e)
    }
  }

  function remove(key: StorageKeyType) {
    try {
      storage.removeItem(key)
    } catch (e) {
      console.error(e)
    }
  }

  function clear() {
    try {
      storage.clear()
    } catch (e) {
      console.error(e)
    }
  }

  function key<T>(index: number, options?: Options<T>): T | null {
    try {
      const raw = storage.key(index)
      if (raw) {
        return deserializer(raw, options)
      }
    } catch (e) {
      console.error(e)
    }
    if (isFunction(options?.defaultValue)) {
      return options?.defaultValue() ?? null
    }
    return options?.defaultValue ?? null
  }

  return {
    get,
    set,
    remove,
    clear,
    key,
    length: storage.length,
  }
}
