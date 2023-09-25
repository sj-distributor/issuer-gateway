/* eslint-disable no-var */
declare type ToastHandle = import("./components/toast").ToastHandle

declare module globalThis {
  var $toast: ToastHandle
}
